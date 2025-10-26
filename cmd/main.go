package main

import (
	"chitchat/pkg/gateway"
	"chitchat/pkg/protocol"
)

func main() {
	newServer := gateway.NewServer(protocol.Gossip{})

	newServer.Serve()
}
