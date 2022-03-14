package main

import (
	"net"
)

type Server interface {
	handleConn(conn net.Conn, ch chan<- error)
	Start() error
	validateResult(rightAnswer, result interface{}) bool
}
