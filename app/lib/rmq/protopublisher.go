package rmq

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type ProtoPublisher struct {
	Conn *Connection
}

func NewProtoPublisher(conn *Connection) *ProtoPublisher {
	return &ProtoPublisher{Conn: conn}
}

func (p *ProtoPublisher) Publish(exchange, routingKey string, message proto.Message, publish *amqp.Publishing) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal the message to protobuf")
	}

	publish.Body = data
	if err = p.Conn.Publish(exchange, routingKey, publish); err != nil {
		return errors.WithMessage(err, "failed to publish the protobuf message")
	}
	return nil
}
