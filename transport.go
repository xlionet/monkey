package monkey

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Transport 传输层抽象
type Transport interface {
	ReadData() (*Envelope, error)
	SendData(*Envelope) error
	IsClosed() bool
	Close()
	GetTag(interface{}) interface{}
	SetTag(interface{}, interface{})
	// HandlerErrFunc
	HandleError(error)
}

// WSTransport websocket transport
type WSTransport struct {
	// Request      *http.Request
	Keys   map[interface{}]interface{}
	conn   *websocket.Conn
	output chan *Envelope

	// monkey  *Monkey
	open     bool
	rwmutex  *sync.RWMutex
	quitSig  chan int
	wg       *sync.WaitGroup
	innerErr error

	conf *TransportConfig

	//TODO processor
}

var (
	// ErrTransportClosed ...
	ErrTransportClosed = errors.New("transport already closed")
	// ErrReadChanClosed ...
	ErrReadChanClosed = errors.New("read chan cloesed")
	// ErrSendBufferFull 发送的缓存队列已经满了
	ErrSendBufferFull = errors.New("send chan buffer is full")
)

// NewWSTransport new a ws tranport instance
func NewWSTransport(conn *websocket.Conn, keys map[interface{}]interface{}) (ts *WSTransport) {
	ts = &WSTransport{
		Keys: keys,
		conn: conn,

		rwmutex: &sync.RWMutex{},
		quitSig: make(chan int),
		wg:      &sync.WaitGroup{},
		conf:    NewTransportConfig(),
	}

	ts.output = make(chan *Envelope, ts.conf.SendBufChanSize)
	return
}

func (s *WSTransport) beginWork(protocol Protocol) {
	s.conn.SetPongHandler(func(appData string) error { return protocol.OnPong([]byte(appData)) })

	s.rwmutex.Lock()
	s.open = true //必须同步调用
	s.rwmutex.Unlock()

	s.wg.Add(1)
	go s.writePump(protocol)
	fmt.Println("monkey start")
}

func (s *WSTransport) writePump(protocol Protocol) {
	ticker := time.NewTicker(time.Duration(s.conf.PingInterval) * time.Second)
	defer func() {
		ticker.Stop()
		s.wg.Done()
	}()

	for s.innerErr == nil {
		select {
		case msg, ok := <-s.output:
			if !ok {
				return
			}

			err := s.writeRaw(msg)
			if err != nil {
				// 本次发送失败，发送了网络问题？需要断开链接 触发重连
				s.innerErr = err
				s.HandleError(err)
				return
			}

			// 收到客户端的退出信息
			if msg.T == websocket.CloseMessage {
				fmt.Println("exit: recv close message")
				return
			}

		case <-ticker.C:
			err := s.ping(protocol.OnPing(nil))
			if err != nil {
				s.HandleError(err)
			}

		case <-s.quitSig:
			//可能会有没有读完的 buff 包 TODO: 需要处理掉，不然会丢包

			fmt.Println("exit: recv quit signal")
			return
		}
	}
}

// ReadData ...
func (s *WSTransport) ReadData() (*Envelope, error) {
	if s.innerErr != nil {
		return nil, s.innerErr
	}

	t, msg, err := s.conn.ReadMessage()
	if err != nil {
		s.HandleError(err)
		return nil, err
	}
	return &Envelope{T: t, Msg: msg}, nil

}

func (s *WSTransport) close() {
	if s.closed() {
		panic("WSTransport already closed")
	}

	// send close msg to close conn
	err := s.writeRaw(&Envelope{T: websocket.CloseMessage, Msg: []byte{}})
	if err != nil {
		fmt.Println("send close message to client failed with err:", err)
	}

	close(s.output)
	s.rwmutex.Lock()
	s.open = false
	s.rwmutex.Unlock()
	s.quitSig <- 1

}

func (s *WSTransport) closed() bool {
	s.rwmutex.RLock()
	defer s.rwmutex.RUnlock()

	return !s.open
}

func (s *WSTransport) writeRaw(message *Envelope) error {
	if s.closed() {
		return errors.New("tried to write to a closed session")
	}

	err := s.conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(s.conf.WriteMaxDurtion)))
	if err != nil {
		return err
	}

	err = s.conn.WriteMessage(message.T, message.Msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *WSTransport) ping(msg []byte) error {
	return s.writeRaw(&Envelope{T: websocket.PingMessage, Msg: msg})
}

// SendData send a envelop msg
func (s *WSTransport) SendData(message *Envelope) error {
	if message == nil {
		return nil
	}

	if s.closed() {
		return ErrTransportClosed
	}

	if s.innerErr != nil {
		return s.innerErr
	}

	select {
	case s.output <- message:
	default:
		s.HandleError(ErrSendBufferFull)
		return ErrSendBufferFull
	}

	return nil
}

// Close close transport
func (s *WSTransport) Close() {
	s.closeByClient()
}

// IsClosed  is closed conn
func (s *WSTransport) IsClosed() bool {
	return s.closed()
}

// SetTag set tag
func (s *WSTransport) SetTag(key interface{}, value interface{}) {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	s.Keys[key] = value
}

// GetTag get tag value
func (s *WSTransport) GetTag(key interface{}) interface{} {
	s.rwmutex.Lock()
	defer s.rwmutex.Unlock()
	return s.Keys[key]
}

// HandleError err 处理
func (s *WSTransport) HandleError(err error) {
	// log(err)
	fmt.Printf("meet error:%v", err)
}

// 退出有两种，1 服务端主动让客户端退出， 2 客户端自己触发的退出，有可能都没有发出退出的信息msg
func (s *WSTransport) closeByServer() { // nolint: unused
	// 1. 先发送退出包，不管是否能发出
	s.SendData(&Envelope{T: websocket.CloseMessage, Msg: []byte{}}) // nolint: errcheck
	// 2. 服务端发出退出信号回收服务端资源，
	s.quitSig <- 1
	s.wg.Wait() // 等待关闭

	// 3. 关闭链接
	s.conn.Close() // 回收

	close(s.output)
	s.rwmutex.Lock()
	s.open = false
	s.rwmutex.Unlock()

}

// 收到退出包，将要退出
func (s *WSTransport) closeByClient() {
	// 2. 服务端发出退出信号回收服务端资源，
	s.quitSig <- 1
	s.wg.Wait() // 等待关闭

	// 3. 关闭链接
	s.conn.Close() // 回收

	close(s.output)
	s.rwmutex.Lock()
	s.open = false
	s.rwmutex.Unlock()
}
