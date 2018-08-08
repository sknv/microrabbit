package server

import (
	"context"
	"log"
	"time"

	"github.com/streadway/amqp"

	"github.com/golang/protobuf/proto"
	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/math/rpc"
)

func RegisterMathServer(rmqServer *rmq.Server, math rpc.Math) {
	mathServer := newMathServer(rmqServer.Conn, math)
	mathServer.route(rmqServer)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type mathServer struct {
	math      rpc.Math
	publisher *rmq.ProtoPublisher
}

func newMathServer(rconn *amqp.Connection, math rpc.Math) *mathServer {
	return &mathServer{
		math:      math,
		publisher: rmq.NewProtoPublisher(rconn),
	}
}

// map a request to a pattern
func (s *mathServer) route(rmqServer *rmq.Server) {
	rmqServer.Handle(rpc.CirclePattern, false, false, 0, withLogger(s.circle))
	// rmqServer.Handle(rpc.RectPattern, false, false, 0, withLogger(s.rect))
}

func (s *mathServer) circle(ctx context.Context, message *amqp.Delivery) {
	args := new(rpc.CircleArgs)
	if err := proto.Unmarshal(message.Body, args); err != nil {
		panic(err) // todo: return error
	}

	reply, err := s.math.Circle(ctx, args)
	if err != nil {
		panic(err) // todo: return error
	}

	if err = s.publisher.Publish("", message.ReplyTo, message.CorrelationId, reply); err != nil {
		panic(err) // todo: return error
	}
}

// func (s *mathServer) rect(ctx context.Context, message *rmq.Msg) {
// 	args := new(rpc.RectArgs)
// 	if err := proto.Unmarshal(message.Data, args); err != nil {
// 		panic(err) // todo: return error
// 	}

// 	reply, err := s.math.Rect(ctx, args)
// 	if err != nil {
// 		panic(err) // todo: return error
// 	}

// 	if err = s.publisher.Publish(message.Reply, reply); err != nil {
// 		panic(err) // todo: return error
// 	}
// }

// ----------------------------------------------------------------------------
// middleware example
// ----------------------------------------------------------------------------

func withLogger(next rmq.HandlerFunc) rmq.HandlerFunc {
	return func(ctx context.Context, msg *amqp.Delivery) {
		start := time.Now()
		defer func() {
			log.Printf("[INFO] request \"%s\" processed in %s", msg.RoutingKey, time.Since(start))
		}()
		next(ctx, msg)
	}
}
