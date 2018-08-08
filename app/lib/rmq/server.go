package rmq

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type ServerProc struct {
	Channel  *Channel
	Messages <-chan amqp.Delivery
}

func (p *ServerProc) Close() {
	if p.Channel != nil {
		p.Channel.Close()
	}
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type Server struct {
	Conn *amqp.Connection

	mu    sync.RWMutex
	procs map[string]*ServerProc
}

func NewServer(conn *amqp.Connection) *Server {
	return &Server{Conn: conn}
}

func (s *Server) Handle(pattern string, durable, autoAck bool) (*ServerProc, error) {
	if pattern == "" {
		return nil, errors.New("rmq: pattern must be present")
	}

	// work with server is thread-safe
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.procs[pattern]; exist {
		return nil, errors.New("rmq: multiple registrations for " + pattern)
	}

	// open a channel per a pattern
	ch, err := NewChannel(s.Conn)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open a channel for the server")
	}

	// declare a queue
	queue, err := ch.DeclareQueue(pattern, durable)
	if err != nil {
		ch.Close() // close the channel opened for this pattern now
		return nil, errors.WithMessage(err, "failed to declare a queue for the server")
	}

	msgs, err := ch.Consume(queue.Name, autoAck)
	if err != nil {
		ch.Close() // close the channel opened for this pattern now
		return nil, errors.WithMessage(err, "failed to register a consumer for the server")
	}

	// if all is ok, store the proc to close it later
	if s.procs == nil {
		s.procs = make(map[string]*ServerProc)
	}
	srvProc := &ServerProc{
		Channel:  ch,
		Messages: msgs,
	}
	s.procs[pattern] = srvProc
	return srvProc, nil
}

func (s *Server) Stop() {
	for _, proc := range s.procs {
		proc.Close()
	}
}
