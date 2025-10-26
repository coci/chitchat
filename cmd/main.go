package main

import (
	"chitchat/internal/gateway"
	"chitchat/internal/protocol"
)

func main() {
	newServer := gateway.NewServer(protocol.Gossip{})

	newServer.Serve()
}
