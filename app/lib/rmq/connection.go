package rmq

import (
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	reconnectTimeout = time.Second
)

type Connection struct {
	*amqp.Connection

	mu sync.Mutex
}

func Dial(addr string) (*Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}
	return &Connection{Connection: conn}, nil
}

func DialWithReconnect(addr string) (*Connection, error) {
	conn, closed, err := dialWithNotify(addr)
	if err != nil {
		return nil, err
	}

	rmqConn := &Connection{Connection: conn}
	rmqConn.ReconnectOnCloseAsync(addr, closed)
	return rmqConn, nil
}

func (c *Connection) OpenChannel() (*Channel, error) {
	ch, err := NewChannel(c.Connection)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (c *Connection) Publish(exchange, routingKey string, message *amqp.Publishing) error {
	ch, err := c.OpenChannel()
	if err != nil {
		return errors.WithMessage(err, "failed to open a channel for publishing")
	}
	defer ch.Close()

	if err = ch.PublishTo(exchange, routingKey, message); err != nil {
		return errors.WithMessage(err, "failed to publish the message")
	}
	return nil
}

func (c *Connection) ReconnectOnCloseAsync(addr string, closed <-chan *amqp.Error) {
	go func() {
		for {
			_, exist := <-closed
			if !exist { // intentional close, do not try to reconnect
				return
			}

			log.Print("[ERROR] lost the connection to RabbitMQ")
			newConn, newClosed := redialWithNotify(addr)
			log.Print("[INFO] reconnected to RabbitMQ")
			closed = newClosed // replace the close channel

			c.mu.Lock()
			c.Connection = newConn // replace the connection
			c.mu.Unlock()
		}
	}()
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func dialWithNotify(addr string) (*amqp.Connection, <-chan *amqp.Error, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, nil, err
	}

	closed := make(chan *amqp.Error)
	conn.NotifyClose(closed)
	return conn, closed, nil
}

func redialWithNotify(addr string) (*amqp.Connection, <-chan *amqp.Error) {
	for {
		conn, closed, err := dialWithNotify(addr)
		if err == nil {
			return conn, closed
		}
		time.Sleep(reconnectTimeout)
	}
}
