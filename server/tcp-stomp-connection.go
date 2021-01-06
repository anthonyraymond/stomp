package server

import (
	"io"
	"net"
	"time"
)

type TcpStompConnection struct {
	conn net.Conn
}

func (c *TcpStompConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *TcpStompConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TcpStompConnection) Send(data []byte) error {
	_, err := c.conn.Write(data)
	return err
}

func (c *TcpStompConnection) Receive() (io.Reader, error) {
	return c.conn, nil
}

func (c *TcpStompConnection) SetReadDeadline(deadline time.Time) error {
	return c.conn.SetReadDeadline(deadline)
}

func (c *TcpStompConnection) Close() error {
	return c.conn.Close()
}
