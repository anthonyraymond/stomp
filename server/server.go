/*
Package server contains a simple STOMP server implementation.
*/
package server

import (
	"time"
)

// The STOMP server has the concept of queues and topics. A message
// sent to a queue destination will be transmitted to the next available
// client that has subscribed. A message sent to a topic will be
// transmitted to all subscribers that are currently subscribed to the
// topic.
//
// Destinations that start with this prefix are considered to be queues.
// Destinations that do not start with this prefix are considered to be topics.
const QueuePrefix = "/queue"

// Default server parameters.
const (
	// Default read timeout for heart-beat.
	// Override by setting Server.HeartBeat.
	DefaultHeartBeat = time.Minute
)

// Interface for authenticating STOMP clients.
type Authenticator interface {
	// Authenticate based on the given login and passcode, either of which might be nil.
	// Returns true if authentication is successful, false otherwise.
	Authenticate(login, passcode string) bool
}

// A Server defines parameters for running a STOMP server.
type Server struct {
	Authenticator Authenticator // Authenticates login/passcodes. If nil no authentication is performed
	QueueStorage  QueueStorage  // Implementation of queue storage. If nil, in-memory queues are used.
	HeartBeat     time.Duration // Preferred value for heart-beat read/write timeout, if zero, then DefaultHeartBeat.
}

func ListenAndServe(entryPoint StompEntryPoint) error {
	s := &Server{}
	return s.ListenAndServe(entryPoint)
}

func (s *Server) ListenAndServe(entryPoint StompEntryPoint) error {
	err := entryPoint.Listen()
	if err != nil {
		return err
	}

	proc := newRequestProcessor(s)
	return proc.Serve(entryPoint)
}
