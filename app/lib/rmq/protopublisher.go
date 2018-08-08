package rmq

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type ProtoPublisher struct {
	Conn *amqp.Connection
}

func NewProtoPublisher(conn *amqp.Connection) *ProtoPublisher {
	return &ProtoPublisher{Conn: conn}
}

func (p *ProtoPublisher) Publish(exchange, routingKey, correlationId string, message proto.Message) error {
	ch, err := NewChannel(p.Conn)
	if err != nil {
		return errors.WithMessage(err, "failed to open a channel for the proto publisher")
	}
	defer ch.Close()

	data, err := proto.Marshal(message)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal a message to protobuf")
	}

	publish := &amqp.Publishing{Body: data}
	if correlationId != "" {
		publish.CorrelationId = correlationId
	}

	if err = ch.Publish(exchange, routingKey, publish); err != nil {
		return errors.WithMessage(err, "failed to publish a proto message")
	}
	return nil
}
