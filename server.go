package main

import (
	"net"
)

type ServerOptions interface {
	SetSaltNumber(saltNumber int)
	SetNonceNumber(nonceNumber int)
	SetMaxRepeatNumber(maxNumber int)
	SetMinRepeatNumber(minNumber int)
}

type Server interface {
	HandleConn(conn net.Conn, ch chan<- error)
	Start() error
	ServerOptions
}
