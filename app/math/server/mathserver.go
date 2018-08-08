package server

// import (
// 	"context"
// 	"log"

// 	"github.com/pkg/errors"
// 	"github.com/streadway/amqp"

// 	"github.com/sknv/microrabbit/app/lib/rmq"
// 	"github.com/sknv/microrabbit/app/math/rpc"
// )

// type MathServer struct {
// 	*rmq.Server

// 	math Math
// 	done chan struct{}
// }

// func NewMathServer(rconn *amqp.Connection) *MathServer {
// 	return &MathServer{
// 		Server: rmq.NewServer(rconn),
// 		math:   Math{},
// 		done:   make(chan struct{}),
// 	}
// }

// func (s *MathServer) ServeAsync() error {
// 	// register method handlers here
// 	//
// 	circleProc, err := s.Handle(rpc.CirclePattern, false, false)
// 	if err != nil {
// 		return errors.WithMessage(err, "failed to consume the circle pattern")
// 	}
// 	rectProc, err := s.Handle(rpc.RectPattern, false, false)
// 	if err != nil {
// 		return errors.WithMessage(err, "failed to consume the rect pattern")
// 	}

// 	// base context
// 	ctx := context.Background()

// 	// listen for messages in a goroutine
// 	//
// 	go func() {
// 		run := true
// 		for run {
// 			select {
// 			case circleMsg := <-circleProc.Messages:
// 				s.handleCircle(ctx, circleProc.Channel, &circleMsg)
// 			case rectMsg := <-rectProc.Messages:
// 				s.handleCircle(ctx, rectProc.Channel, &rectMsg)
// 			case <-s.done:
// 				run = false
// 			}
// 		}
// 	}()

// 	log.Print("[INFO] starting a math server")
// 	return nil
// }

// func (s *MathServer) Stop() {
// 	s.done <- struct{}{}
// 	s.Server.Stop()
// 	log.Print("[INFO] math server stopped")
// }

// // ----------------------------------------------------------------------------
// // ----------------------------------------------------------------------------
// // ----------------------------------------------------------------------------

// func (s *MathServer) handleCircle(ctx context.Context, channel *rmq.Channel, message *amqp.Delivery) {

// }
