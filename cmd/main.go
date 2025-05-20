package main

import (
	"fmt"
	"github.com/coci/chitchat/pkg/protocol"
	"github.com/coci/chitchat/pkg/utils"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		msg, err := protocol.ParseMessage(conn)
		if err != nil {
			fmt.Println("Parse error:", err)
			return
		}
		fmt.Printf("Received: Opcode=0x%02X, Length=%d, Body=%s\n",
			msg.Header.Opcode, msg.Header.Length, string(msg.Body))
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	utils.MakeSplashScreen()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
