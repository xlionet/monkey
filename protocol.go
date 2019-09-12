package monkey

// Protocol ...
type Protocol interface {
	OnTransportMade(Transport)
	OnTransportLost(Transport)
	OnTransportData(Transport, *Envelope)
	OnPing([]byte) []byte
	OnPong([]byte) error
}

// WSProtocol protocol implement for monkey websockt
type WSProtocol struct{}

// OnTransportMade ...
func (p *WSProtocol) OnTransportMade(transport Transport) {

}

// OnTransportLost ...
func (p *WSProtocol) OnTransportLost(transport Transport) {

}

// OnTransportData ...
func (p *WSProtocol) OnTransportData(transport Transport, env *Envelope) {

}

// OnPing ...
func (p *WSProtocol) OnPing(pp []byte) []byte {

	return append(pp, []byte("my ping")...)

}

// OnPong ...
func (p *WSProtocol) OnPong(pp []byte) error {

	return nil
}

// NewWSProtocol ...
func NewWSProtocol() *WSProtocol {
	return &WSProtocol{}
}
