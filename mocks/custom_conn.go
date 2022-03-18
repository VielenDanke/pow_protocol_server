package mocks

import (
	"net"
	"time"
)

type CustomConn struct {
}

func (c CustomConn) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	bytesToRead := []byte("user,string,string,string")
	for i, v := range bytesToRead {
		b[i] = v
	}
	return len(bytesToRead), nil
}

func (c CustomConn) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	bytesToWrite := []byte("user,string,string,string")
	for i, v := range bytesToWrite {
		b[i] = v
	}
	return len(bytesToWrite), nil
}

func (c CustomConn) Close() error {
	return nil
}

func (c CustomConn) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   []byte("10.4.108.11"),
		Port: 8080,
	}
}

func (c CustomConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   []byte("10.4.108.11"),
		Port: 8080,
	}
}

func (c CustomConn) SetDeadline(t time.Time) error {
	return nil
}

func (c CustomConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c CustomConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func NewCustomConn() net.Conn {
	return &CustomConn{}
}
