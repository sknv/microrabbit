package server

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/math/rpc"
)

func RegisterMathServer(server *rmq.Server, math rpc.Math) {
	mathServer := newMathServer(server.Conn, math)
	mathServer.route(server)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type mathServer struct {
	math      rpc.Math
	publisher *rmq.ProtoPublisher
}

func newMathServer(conn *rmq.Connection, math rpc.Math) *mathServer {
	return &mathServer{
		math:      math,
		publisher: rmq.NewProtoPublisher(conn),
	}
}

// map a request to a pattern
func (s *mathServer) route(server *rmq.Server) {
	server.Handle(rpc.CirclePattern, s.circle)
	server.Handle(rpc.RectPattern, s.rect)
}

func (s *mathServer) circle(ctx context.Context, message *amqp.Delivery) error {
	args := new(rpc.CircleArgs)
	if err := proto.Unmarshal(message.Body, args); err != nil {
		return err
	}

	reply, err := s.math.Circle(ctx, args)
	if err != nil {
		return err
	}

	publish := &amqp.Publishing{CorrelationId: message.CorrelationId}
	if err = s.publisher.Publish("", message.ReplyTo, reply, publish); err != nil {
		return err
	}
	return nil
}

func (s *mathServer) rect(ctx context.Context, message *amqp.Delivery) error {
	args := new(rpc.RectArgs)
	if err := proto.Unmarshal(message.Body, args); err != nil {
		return err
	}

	reply, err := s.math.Rect(ctx, args)
	if err != nil {
		return err
	}

	publish := &amqp.Publishing{CorrelationId: message.CorrelationId}
	if err = s.publisher.Publish("", message.ReplyTo, reply, publish); err != nil {
		return err
	}
	return nil
}
