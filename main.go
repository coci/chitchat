package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func frameToString(f Frame) string {
	return fmt.Sprintf(
		"GOSSIP HEADER {\n"+
			"  Magic: 0x%04X\n"+
			"  Version: %d\n"+
			"  MsgType: %d\n"+
			"  StreamID: %d\n"+
			"  Length: %d\n"+
			"  Payload: %q\n"+
			"}",
		f.Header.Magic,
		f.Header.Version,
		f.Header.MessageType,
		f.Header.StreamId,
		f.Header.Length,
		f.Body,
	)
}

func HandleConnection(conn net.Conn) {
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

		frame, err := EncodeMessage(buff[:n])
		if err != nil {
			log.Println("Decode error:", err)
			continue
		}

		fmt.Println(frameToString(frame))
		break // stop after one frame
	}
}
func main() {
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
		HandleConnection(conn)
	}
}
