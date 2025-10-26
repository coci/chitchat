package gateway

import (
	"chitchat/internal/protocol"
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	protocol protocol.IProtocol
}

func NewServer(protocol protocol.IProtocol) *Server {
	return &Server{
		protocol: protocol,
	}
}

func (s *Server) Serve() {
	fmt.Println("ChitChat server starting ....")

	tcp, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	defer tcp.Close()

	for {
		conn, err := tcp.Accept()
		if err != nil {
			log.Fatal(err)
		}
		s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Welcome to ChitChat\n"))

	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed connection.")
				break
			}
			log.Println("Read error:", err)
			break
		}

		if n == 0 {
			// connection open but no data
			break
		}

		frame, err := s.protocol.ParseMessage(buff[:n])

		if err != nil {
			log.Println("Decode error:", err)
			continue
		}

		fmt.Println(s.protocol.FrameToString(frame))
	}
}
