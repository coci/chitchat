package main

import (
	"github.com/coci/chitchat/pkg/protocol"
	"github.com/coci/chitchat/pkg/utils"
	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		msg, err := protocol.ParseMessage(conn)
		if err != nil {
			log.Fatal("Error parsing message: ", err)
			return
		}

		switch msg.Header.Opcode {
		case protocol.CreateUser:
			// handle create user
		case protocol.GetUserToken:
			// handle get user token
		case protocol.GetChatroomLists:
			// handle get chatroom lists
		case protocol.CreateChatroom:
			// handle create chatroom
		case protocol.JoinChatroom:
			// handle join chatroom
		case protocol.GetChatroomUser:
			// handle get chatroom user
		case protocol.GetChatroomMessages:
			// handle get chatroom messages
		case protocol.SendMessageTOChatroom:
			// handle send message to chatroom
		case protocol.UserLoggedOut:
			// handle user logged out
		default:
			// handle unknown opcode
		}
	}
}

func main() {
	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":9000" // default fallback
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	utils.MakeSplashScreen()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}
