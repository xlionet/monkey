package monkey

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// type handleMessageFunc func(*WSTransport, []byte)
type handleErrorFunc func(*WSTransport, error)

// Monkey ...
type Monkey struct {
	Upgrader *websocket.Upgrader
	// errorHandler handleErrorFunc
	WsProtocol Protocol
	Open       bool
	conf       *Config
}

// ErrMonkeyIsNotOpen ...
var ErrMonkeyIsNotOpen = errors.New("monkey instance is not open")

// NewMonkey new
func New(pt Protocol, config *Config) *Monkey {
	m := &Monkey{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:   1024,
			HandshakeTimeout: time.Second * 10,
			WriteBufferSize:  1024,
			CheckOrigin:      func(r *http.Request) bool { return true },
		},

		WsProtocol: pt,
		conf:       config,
		Open:       true,
	}
	return m
}

// HandleConnection Monkey upgrades http requests to websocket connections
func (m *Monkey) HandleConnection(w http.ResponseWriter, r *http.Request) error {
	return m.HandleConnectionWithKeys(w, r, nil)
}

// HandleConnectionWithKeys does the same as HandleRequest but populates monkey.Keys with keys.
func (m *Monkey) HandleConnectionWithKeys(w http.ResponseWriter, r *http.Request, keys map[interface{}]interface{}) error {
	if !m.Open {
		return ErrMonkeyIsNotOpen
	}

	conn, err := m.Upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		return err
	}

	ts := NewWSTransport(conn, keys)
	ts.beginWork(m.WsProtocol)

	m.WsProtocol.OnTransportMade(ts)
	m.startProcess(ts, m.WsProtocol)
	defer func() {
		ts.close()
		m.WsProtocol.OnTransportLost(ts)
	}()
	return nil
}

func (m *Monkey) startProcess(transport Transport, protocol Protocol) {

	for m.Open && !transport.IsClosed() {
		data, err := transport.ReadData()
		if err != nil {
			transport.HandleError(err)
			return
		}

		protocol.OnTransportData(transport, data)
	}
}

// Serve start server
func (m *Monkey) Serve(handler http.Handler) error {
	addr := fmt.Sprintf("0.0.0.0:%d", m.conf.WSListernPort)
	return http.ListenAndServe(addr, handler)
}
