package server

import (
	"github.com/gorilla/websocket"
	"io"
	"net"
	"time"
)

type WebSocketStompConnection struct {
	conn *websocket.Conn
}

func (c *WebSocketStompConnection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *WebSocketStompConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *WebSocketStompConnection) Send(data []byte) error {
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (c *WebSocketStompConnection) Receive() (io.Reader, error) {
	_, r, err := c.conn.NextReader()
	return r, err
}

func (c *WebSocketStompConnection) SetReadDeadline(deadline time.Time) error {
	return c.conn.SetReadDeadline(deadline)
}

func (c *WebSocketStompConnection) Close() error {
	return c.conn.Close()
}
