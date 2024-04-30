package p2p

import "net"

// RCP holds any arbitrary data
// that is sent over each transport between two nodes
type RPC struct {
	Payload []byte
	From    net.Addr
}
