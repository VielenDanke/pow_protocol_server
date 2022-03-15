package main

import (
	"net"
)

type Server interface {
	HandleConn(conn net.Conn, ch chan<- error)
	Start() error
	SetSaltNumber(saltNumber int)
	SetNonceNumber(nonceNumber int)
	SetMaxRepeatNumber(maxNumber int)
	SetMinRepeatNumber(minNumber int)
}
