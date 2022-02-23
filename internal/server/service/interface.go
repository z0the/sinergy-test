package service

import (
	"net"
)

type Service interface {
	ServeConn(conn net.Conn)
}
