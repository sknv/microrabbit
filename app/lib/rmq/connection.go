package rmq

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/xcontext"
	"github.com/sknv/microrabbit/app/lib/xstrings"
)

const (
	directReplyQueue = "amq.rabbitmq.reply-to"
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

func (c *Connection) Request(ctx context.Context, routingKey string, message *amqp.Publishing) (*amqp.Delivery, error) {
	ch, err := c.OpenChannel()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to open a sync channel")
	}
	defer ch.Close()

	// prepare the message
	correlationID := xstrings.RandomString(32)
	msg := prepareMessage(ctx, message, correlationID)

	// consume from the reply queue
	msgs, err := ch.ConsumeFrom(msg.ReplyTo, true)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to register a sync consumer")
	}

	// publish the message
	if err = ch.PublishTo("", routingKey, msg); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("failed to publish the sync message to %s", routingKey))
	}
	return handleReply(ctx, msgs, correlationID) // wait for the reply
}

func (c *Connection) ReconnectOnCloseAsync(addr string, closed <-chan *amqp.Error) {
	go func() {
		for {
			_, exist := <-closed
			if !exist { // intentional close, do not try to reconnect
				return
			}

			newConn, newClosed := redialWithNotify(addr)
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
			log.Print("[INFO] reconnected to RabbitMQ")
			return conn, closed
		}

		log.Print("[ERROR] failed to reconnect to RabbitMQ: ", err)
		time.Sleep(reconnectTimeout)
	}
}

func prepareMessage(ctx context.Context, message *amqp.Publishing, correlationID string) *amqp.Publishing {
	message.ReplyTo = directReplyQueue
	message.CorrelationId = correlationID
	message = expireIfNeeded(ctx, message) // expire a message if a deadline specified for the context
	return message
}

func expireIfNeeded(ctx context.Context, message *amqp.Publishing) *amqp.Publishing {
	timeout, exist := xcontext.Timeout(ctx)
	if !exist {
		return message
	}

	exp := math.Ceil(timeout.Seconds() * 1000) // expiration in ms rounded to the nearest greater int
	message.Expiration = fmt.Sprint(exp)
	return message
}

func handleReply(ctx context.Context, messages <-chan amqp.Delivery, correlationID string) (*amqp.Delivery, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, errors.WithMessage(ctx.Err(), "failed to wait for the reply")
		case msg := <-messages:
			if correlationID == msg.CorrelationId {
				return &msg, nil
			}
		}
	}
}
