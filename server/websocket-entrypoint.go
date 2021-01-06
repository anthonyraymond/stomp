package server

import (
	"context"
	"fmt"
	"github.com/go-stomp/stomp/server/client"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type webSocketEntryPoint struct {
	httpServer      *http.Server
	Port            int
	clientConnected chan client.StompConnection
}

func NewWebSocketEntryPoint(port int) StompEntryPoint {
	return &webSocketEntryPoint{
		httpServer:      nil,
		Port:            port,
		clientConnected: make(chan client.StompConnection),
	}
}

func (w *webSocketEntryPoint) Listen() error {
	handler := http.NewServeMux()
	handler.HandleFunc("/ws", func(writer http.ResponseWriter, r *http.Request) {
		respHeader := http.Header{}
		if h := r.Header.Get("Sec-WebSocket-Protocol"); h != "" {
			respHeader.Add("Sec-WebSocket-Protocol", strings.Split(h, ",")[0])
		}
		conn, err := upgrader.Upgrade(writer, r, respHeader)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		select {
		case w.clientConnected <- &WebSocketStompConnection{conn: conn}:
		case <-time.After(10 * time.Second):
		}
	})

	w.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", w.Port),
		Handler:           handler,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
	}

	go func() {
		err := w.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("http server has failed to start", err)
		}
	}()

	return nil
}

func (w *webSocketEntryPoint) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := w.httpServer.Shutdown(ctx)
	if err != nil {
		_ = w.httpServer.Close()
	}
}

func (w *webSocketEntryPoint) Accept() (client.StompConnection, error) {
	return <-w.clientConnected, nil
}
