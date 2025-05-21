package main

import (
	"context"
	"github.com/coci/chitchat/pkg/protocol"
	"github.com/coci/chitchat/pkg/utils"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func handleConnection(cancelContext context.Context, conn net.Conn) {
	defer wg.Done()
	defer conn.Close()

	for {
		_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		select {
		case <-cancelContext.Done():
			log.Println("Connection canceled, closing handler.")
			return
		default:
			msg, err := protocol.ParseMessage(conn)

			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Println("Error parsing message:", err)
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
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":9000"
	}
	utils.MakeSplashScreen()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					log.Print("Error accepting connection: ", err)
					continue
				}
				wg.Add(1)
				go handleConnection(ctx, conn)
			}
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down the server .....")
	wg.Wait()
	listener.Close()
	log.Println("server shut down")
}
