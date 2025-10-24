package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Welcome to ChitChat\n"))

	for {
		buff := make([]byte, 1024)

		n, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			log.Print(err)
		}
		frame, err := EncodeMessage(buff[:n])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(frame)
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
