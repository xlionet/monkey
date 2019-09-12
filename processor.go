package monkey

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrMsgTypeAlreadRegisted ...
	ErrMsgTypeAlreadRegisted = errors.New("msg type has alread registed")
)

// HandleFunc 包处理器
type HandleFunc func(context.Context, *JSONGatePacket, Transport)

// JSONPacketProcess 。。。
type JSONPacketProcess struct {
	dispatchFunc map[uint16]HandleFunc
	timeoutFunc  map[uint16]HandleFunc
	sync.RWMutex
	Packeter GatePacketer
}

// RegisterFunc 注册 handle
func (jp *JSONPacketProcess) RegisterFunc(msgType uint16, handFunc, timeoutFunc HandleFunc) error {
	return jp.registerFunc(msgType, handFunc, timeoutFunc)
}

func (jp *JSONPacketProcess) registerFunc(msgType uint16, handFunc, timeoutFunc HandleFunc) error {
	jp.Lock()
	defer jp.Unlock()
	if jp.dispatchFunc == nil {
		jp.dispatchFunc = make(map[uint16]HandleFunc)
		jp.timeoutFunc = make(map[uint16]HandleFunc)
	}
	if _, ok := jp.dispatchFunc[msgType]; ok {
		return ErrMsgTypeAlreadRegisted
	}

	jp.dispatchFunc[msgType] = handFunc
	jp.timeoutFunc[msgType] = timeoutFunc

	return nil
}

// GetHandler 获取handle
func (jp *JSONPacketProcess) GetHandler(packetID uint16) (process HandleFunc, timeout HandleFunc) {
	jp.RLock()
	defer jp.RUnlock()

	if handleFunc, ok := jp.dispatchFunc[packetID]; ok {
		process = handleFunc
	}

	if to, ok := jp.timeoutFunc[packetID]; ok {
		timeout = to
	}

	return
}

// OnTransportMade implement for protocol
func (jp *JSONPacketProcess) OnTransportMade(trasnport Transport) {
	// Do nothing but 可以做用户登陆的信息存储
}

// OnTransportLost implement for protocol
func (jp *JSONPacketProcess) OnTransportLost(trasnport Transport) {
	// DO nothing
}

// OnTransportData implement for protocol
func (jp *JSONPacketProcess) OnTransportData(transport Transport, envelop *Envelope) {
	if envelop == nil {
		return
	}

	var base JSONGatePacket
	err := jp.Packeter.Unpack(envelop.Msg, &base)
	if err != nil {
		fmt.Println("unpack failed", err)
		return
	}

	p, _ := jp.GetHandler(uint16(base.PacketID)) // 超时的处理器暂不实现,主要是想针对服务端chan buffer 满的情况，提供给上层处理
	if p == nil {
		return
	}

	// error 不需要抛到这一层，因为这是业务逻辑代码，出错了这里也处理不了
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("process message encounter panic:", err)
		}
	}()
	ctx := context.Background()
	p(ctx, &base, transport)
}

// OnPing implement for protocol
func (jp *JSONPacketProcess) OnPing(m []byte) []byte {
	fmt.Println("on ping")
	return m
}

// OnPong implement for protocol
func (jp *JSONPacketProcess) OnPong(_ []byte) error {
	return nil
}
