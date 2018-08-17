package rmq

import (
	"context"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type HandlerFunc func(context.Context, *amqp.Delivery) error

type ServerInterceptor func(HandlerFunc) HandlerFunc

type Server struct {
	Conn         *Connection
	Interceptors []ServerInterceptor

	mu      sync.Mutex
	entries map[string]*serverEntry
}

func NewServer(conn *Connection, interceptors ...ServerInterceptor) *Server {
	return &Server{
		Conn:         conn,
		Interceptors: interceptors,
	}
}

func (s *Server) Handle(pattern string, handlerFn HandlerFunc) {
	if pattern == "" {
		log.Fatal("rmq: invalid pattern")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.entries[pattern]; exist {
		log.Fatal("rmq: multiple registrations for " + pattern)
	}

	if s.entries == nil {
		s.entries = make(map[string]*serverEntry)
	}
	s.entries[pattern] = newServerEntry(s.Conn, pattern, handlerFn, s.Interceptors)
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
	channel      *Channel
	handlerFn    HandlerFunc
	interceptors []ServerInterceptor
	msgs         <-chan amqp.Delivery
	closed       chan *amqp.Error
}

func newServerEntry(conn *Connection, pattern string, handlerFn HandlerFunc, interceptors []ServerInterceptor) *serverEntry {
	ch, err := conn.OpenChannel()
	if err != nil {
		log.Fatal("rmq: failed to open a channel for the server entry")
	}

	// prevent receiving error messages on channel close
	closed := make(chan *amqp.Error)
	ch.NotifyClose(closed)

	queue, err := ch.DeclareQueue(pattern, false)
	if err != nil {
		log.Fatal("rmq: failed to declare a queue for the server entry")
	}

	msgs, err := ch.ConsumeFrom(queue.Name, true)
	if err != nil {
		log.Fatal("rmq: failed to consume a queue for the server entry")
	}

	return &serverEntry{
		channel:      ch,
		handlerFn:    handlerFn,
		interceptors: interceptors,
		msgs:         msgs,
		closed:       closed,
	}
}

func (e *serverEntry) serveAsync() {
	go func() {
		run := true
		for run {
			select {
			case msg := <-e.msgs:
				e.handleMessageAsync(&msg)
			case <-e.closed:
				run = false
			}
		}
	}()
}

func (e *serverEntry) stop() {
	if e.channel != nil {
		e.channel.Close() // will also notify "closed" channel
	}
}

func (e *serverEntry) handleMessageAsync(message *amqp.Delivery) {
	go func() { // process messages in a goroutine
		// recover on panic
		defer func() {
			if rvr := recover(); rvr != nil {
				log.Printf("[PANIC] recover: %s", rvr)
			}
		}()

		// fill the context with the metadata
		ctx := context.Background()
		metadata := Metadata(message)
		if metadata != nil {
			for key, val := range metadata {
				ctx = ContextWithMetaValue(ctx, key, val)
			}
		}

		// execute the interceptors
		handlerFn := e.chainInterceptors(e.handlerFn)

		// call the handlerFn func
		if err := handlerFn(ctx, message); err != nil {
			log.Printf("[ERROR] failed to handle the message %s: %s", message.RoutingKey, err)
		}
	}()
}

func (e *serverEntry) chainInterceptors(endpoint HandlerFunc) HandlerFunc {
	if len(e.interceptors) == 0 {
		return endpoint
	}

	handler := e.interceptors[len(e.interceptors)-1](endpoint)
	for i := len(e.interceptors) - 2; i >= 0; i-- {
		handler = e.interceptors[i](handler)
	}
	return handler
}
