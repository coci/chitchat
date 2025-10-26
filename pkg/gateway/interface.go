package gateway

import "net"

type IGateway interface {
	Serve()
	HandleConnection(conn net.Conn)
}
