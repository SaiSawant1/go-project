package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represent the remote node over a TCP established connection
type TCPPeer struct {
	// conn is the underlying connection of the peers
	conn net.Conn
	// if we dial and retrive connection = true
	// if we accept and retrive connection = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implement Peer interface.
func (p *TCPPeer) Close() error {
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}

type TCPTrasportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTrasportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTrasportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTrasportOpts: opts,
		rpcch:           make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	return nil
}

// Consume implements the transport interface, which will
// return read-only channel for reading the channel
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP error: %v\n", err)
		}

		fmt.Printf("new incoming connection %+v\n", conn)
		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("Dropping peer connection,%+v\n", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, true)
	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}
