package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

var (
	ProtocolName    = "X"
	ProtocolVersion = 1
	ProtocolLength  = 1
	bootnodeAddr    = "enode://abc...@127.0.0.1:30310"
)

func bootnodeDisco(nodelist []string) []*discover.Node {
	var nodes []*discover.Node
	for _, url := range nodelist {
		if url == "" {
			continue
		}
		node, err := discover.ParseNode(url)
		if err != nil {
			log.Println(fmt.Sprintf("Node URL %s: %v\n", url, err))
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func main() {
	nodeKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	srv := p2p.Server{
		Config: p2p.Config{
			MaxPeers:       10,
			PrivateKey:     nodeKey,
			Name:           "x-chain",
			ListenAddr:     ":30303",
			BootstrapNodes: bootnodeDisco([]string{bootnodeAddr}),
			Protocols:      []p2p.Protocol{ProtocolX()},
		},
	}
	srv.Start()
	defer srv.Stop()

	log.Println("server started")
	fmt.Scanln()
}

func ProtocolX() p2p.Protocol {
	return p2p.Protocol{
		Name:    ProtocolName,
		Version: uint(ProtocolVersion),
		Length:  uint64(ProtocolLength),
		Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
			log.Println("waiting for peer")
			for {
				msg, err := rw.ReadMsg()
				if err != nil {
					return err
				}
				defer msg.Discard()
				log.Println(msg.Code)
				// p2p.Send(rw, 0, newMsg)
			}
		},
	}
}
