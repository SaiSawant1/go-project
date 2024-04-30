package main

import (
	"fmt"
	"log"

	"github.com/SaiSawant1/filestore/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("Doing some logic with the peer")
	return nil
}

func main() {
	TcpOpts := p2p.TCPTrasportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(TcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalln(err)
	}
	select {}
}
