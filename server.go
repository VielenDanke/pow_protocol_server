package main

import (
	"context"
	"net"
)

type CustomConn interface {
	net.Conn
}

type ServerOptions interface {
	SetSaltNumber(saltNumber int)
	SetNonceNumber(nonceNumber int)
	SetMaxRepeatNumber(maxNumber int)
	SetMinRepeatNumber(minNumber int)
	GetMaxRepeatNumber() int
	GetNonceNumber() int
	GetSaltNumber() int
	GetMinRepeatNumber() int
	GetNetworkType() string
	GetAddress() string
}

type Server interface {
	HandleConn(conn net.Conn, ch chan<- error)
	Start(ctx context.Context) error
	ServerOptions
}
