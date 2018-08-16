package rpc

import (
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/status"
)

type ProtoResponder struct {
	*rmq.ProtoPublisher
}

func NewProtoResponder(conn *rmq.Connection) *ProtoResponder {
	return &ProtoResponder{ProtoPublisher: rmq.NewProtoPublisher(conn)}
}

func (r *ProtoResponder) Reply(request *amqp.Delivery, reply proto.Message, err error) error {
	publish := &amqp.Publishing{CorrelationId: request.CorrelationId}
	if err == nil {
		return r.Publish("", request.ReplyTo, publish, reply)
	}

	// transfer error if such exist
	publish = rmq.WithError(publish)
	status, _ := status.FromError(err)
	return r.Publish("", request.ReplyTo, publish, status)
}
