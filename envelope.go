package monkey

// Envelope a msg for send to ws
type Envelope struct {
	T   int
	Msg []byte
	// filter filterFunc
}

// type filterFunc func(*Transport) bool
