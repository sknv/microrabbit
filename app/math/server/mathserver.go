package server

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/rpc"
	"github.com/sknv/microrabbit/app/lib/rmq/status"
	math "github.com/sknv/microrabbit/app/math/rpc"
)

func RegisterMathServer(server *rmq.Server, math math.Math) {
	mathServer := newMathServer(server.Conn, math)
	mathServer.route(server)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type mathServer struct {
	math      math.Math
	responder *rpc.ProtoResponder
}

func newMathServer(conn *rmq.Connection, math math.Math) *mathServer {
	return &mathServer{
		math:      math,
		responder: rpc.NewProtoResponder(conn),
	}
}

// map a request to a pattern
func (s *mathServer) route(server *rmq.Server) {
	server.Handle(math.CirclePattern, s.circle)
	server.Handle(math.RectPattern, s.rect)
}

func (s *mathServer) circle(ctx context.Context, message *amqp.Delivery) error {
	args := new(math.CircleArgs)
	if err := s.decodeArgs(message, args); err != nil {
		return err
	}

	reply, err := s.math.Circle(ctx, args)
	if err = s.responder.Reply(message, reply, err); err != nil {
		return errors.WithMessage(err, "failed to publish the reply")
	}
	return nil
}

func (s *mathServer) rect(ctx context.Context, message *amqp.Delivery) error {
	args := new(math.RectArgs)
	if err := s.decodeArgs(message, args); err != nil {
		return err
	}

	reply, err := s.math.Rect(ctx, args)
	if err = s.responder.Reply(message, reply, err); err != nil {
		return errors.WithMessage(err, "failed to publish the reply")
	}
	return nil
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (s *mathServer) decodeArgs(message *amqp.Delivery, args proto.Message) error {
	err := proto.Unmarshal(message.Body, args)
	if err == nil {
		return nil
	}

	err = errors.WithMessage(err, "failed to unmarshal the message body from protobuf")
	status := status.Error(status.InvalidArgument, err.Error())
	if reperr := s.responder.Reply(message, nil, status); reperr != nil {
		reperr = errors.WithMessage(reperr, "failed to publish the error status")
		err = errors.WithMessage(err, reperr.Error())
	}
	return err
}
