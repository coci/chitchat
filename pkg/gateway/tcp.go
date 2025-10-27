package gateway

import (
	"chitchat/pkg/logger"
	"chitchat/pkg/protocol"
	"io"
	"net"
)

type Server struct {
	IGateway
	protocol protocol.IProtocol
	logger   logger.ILogger
}

func NewServer(protocol protocol.IProtocol, logger logger.ILogger) *Server {
	return &Server{
		protocol: protocol,
		logger:   logger,
	}
}

func (s *Server) Serve() {
	s.logger.Info("Starting server...")

	tcp, err := net.Listen("tcp", ":8080")
	if err != nil {
		s.logger.Error("Failed to Listen", logger.Field{Key: "error", Value: err})
	}

	defer func(tcp net.Listener) {
		err := tcp.Close()
		if err != nil {
			s.logger.Error("Failed to close TCP listener", logger.Field{Key: "error", Value: err})
		}
	}(tcp)

	for {
		conn, err := tcp.Accept()
		if err != nil {
			s.logger.Error("Failed to accept new tcp connection", logger.Field{Key: "error", Value: err})
		}
		s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			s.logger.Error("Failed to close TCP connection", logger.Field{Key: "error", Value: err})
		}
	}(conn)

	_, err := conn.Write([]byte("Welcome to ChitChat\n"))
	if err != nil {
		return
	}

	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				s.logger.Error("Close Connection")
				break
			}
			s.logger.Error("Failed to read data in tcp connection", logger.Field{Key: "error", Value: err})
			break
		}

		if n == 0 {
			s.logger.Error("Close Connection")
			break
		}

		frame, err := s.protocol.ParseMessage(buff[:n])

		if err != nil {
			s.logger.Error("Failed to parse message", logger.Field{Key: "error", Value: err})
			continue
		}

		s.logger.Info(s.protocol.String(frame))
	}
}
