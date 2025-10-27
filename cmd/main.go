package main

import (
	"chitchat/pkg/gateway"
	"chitchat/pkg/logger"
	"chitchat/pkg/protocol"
)

func main() {
	newServer := gateway.NewServer(protocol.Gossip{Logger: logger.NewZapLogger()}, logger.NewZapLogger())

	newServer.Serve()
}
