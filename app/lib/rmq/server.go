package rmq

import (
	"context"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type HandlerFunc func(context.Context, *amqp.Delivery)

type Server struct {
	Conn *amqp.Connection

	mu      sync.RWMutex
	entries map[string]*serverEntry
}

func NewServer(conn *amqp.Connection) *Server {
	return &Server{Conn: conn}
}

func (s *Server) Handle(pattern string, autoAck, durable bool, channelLimit int, handler HandlerFunc) {
	if pattern == "" {
		panic("rmq: invalid pattern")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.entries[pattern]; exist {
		panic("rmq: multiple registrations for " + pattern)
	}

	if s.entries == nil {
		s.entries = make(map[string]*serverEntry)
	}
	s.entries[pattern] = newServerEntry(s.Conn, pattern, autoAck, durable, channelLimit, handler)
}

func (s *Server) ServeAsync() {
	for _, entry := range s.entries {
		entry.serveAsync()
	}
	log.Print("[INFO] rmq server started")
}

func (s *Server) Stop() {
	for _, entry := range s.entries {
		entry.stop()
	}
	log.Print("[INFO] rmq server stopped")
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type serverEntry struct {
	channel *Channel
	autoAck bool
	handler HandlerFunc
	msgs    <-chan amqp.Delivery
	done    chan struct{}
}

func newServerEntry(
	conn *amqp.Connection, pattern string, autoAck, durable bool, channelLimit int, handler HandlerFunc,
) *serverEntry {
	ch, err := NewChannel(conn)
	if err != nil {
		panic("rmq: failed to open a channel for the server entry")
	}

	queue, err := ch.DeclareQueue(pattern, durable)
	if err != nil {
		panic("rmq: failed to declare a queue for the server entry")
	}

	if err = ch.QoS(channelLimit); err != nil {
		panic("rmq: failed to set QoS for the server entry")
	}

	msgs, err := ch.Consume(queue.Name, autoAck)
	if err != nil {
		panic("rmq: failed to consume a queue for the server entry")
	}

	return &serverEntry{
		channel: ch,
		autoAck: autoAck,
		handler: handler,
		msgs:    msgs,
		done:    make(chan struct{}),
	}
}

func (e *serverEntry) serveAsync() {
	go func() {
		run := true
		for run {
			select {
			case msg := <-e.msgs:
				go e.handleMessage(&msg)
			case <-e.done:
				run = false
			}
		}
	}()
}

func (e *serverEntry) handleMessage(message *amqp.Delivery) {
	e.handler(context.Background(), message)
	if !e.autoAck {
		message.Ack(false)
	}
}

func (e *serverEntry) stop() {
	e.done <- struct{}{}
	if e.channel != nil {
		e.channel.Close()
	}
}
