package monkey

// TransportGroup save all of container connections, design for sent to client directly, OnConnection Made to Add OnConnection To Remove
// use Get() func to get a transport to send msg
type TransportGroup interface {
	Add()
	Remove()
	Get()
	GetAll()
}
