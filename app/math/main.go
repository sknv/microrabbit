package main

import (
	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/rmq/interceptors"
	"github.com/sknv/microrabbit/app/lib/xos"
	"github.com/sknv/microrabbit/app/math/cfg"
	"github.com/sknv/microrabbit/app/math/server"
)

func main() {
	cfg := cfg.Parse()

	// connect to RabbitMQ
	consumerConn, err := rmq.DialWithReconnect(cfg.RabbitAddr)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer consumerConn.Close()

	publisherConn, err := rmq.DialWithReconnect(cfg.RabbitAddr)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer publisherConn.Close()

	// handle rmq requests
	srv := rmq.NewServer(consumerConn, interceptors.WithLogger)
	server.RegisterMathServer(publisherConn, srv, &server.MathImpl{})

	// start the rmq server and schedule a stop
	srv.ServeAsync()
	defer srv.Stop()

	// wait for a program exit to stop the rmq server
	xos.WaitForExit()
}
