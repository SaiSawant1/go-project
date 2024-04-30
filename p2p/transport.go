package p2p

// Peer is an interface that represents the remote node.
type Peer interface {
	Close() error
}

// Transport is anything that handels the communication
// between two nodes in a network. This can be of the
// form (TCP,UDP,WebSocket,...).
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
