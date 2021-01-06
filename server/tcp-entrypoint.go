package server

import (
	"github.com/go-stomp/stomp/server/client"
	"net"
)

type tcpEntryPoint struct {
	addr     string
	listener net.Listener
}

func NewTcpEntryPoint(addr string) StompEntryPoint {
	if addr == "" {
		addr = ":61613"
	}
	return &tcpEntryPoint{
		addr: addr,
	}
}

func (t *tcpEntryPoint) Listen() error {
	listener, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	t.listener = listener
	return nil
}

func (t *tcpEntryPoint) Shutdown() {
	if t.listener != nil {
		_ = t.listener.Close()
	}
}

func (t *tcpEntryPoint) Accept() (client.StompConnection, error) {
	c, err := t.listener.Accept()
	if err != nil {
		return nil, err
	}
	return &TcpStompConnection{conn: c}, nil
}
