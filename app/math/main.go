package main

import (
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/xos"
	"github.com/sknv/microrabbit/app/math/cfg"
	"github.com/sknv/microrabbit/app/math/server"
)

func main() {
	cfg := cfg.Parse()

	// connect to RabbitMQ
	rconn, err := amqp.Dial(cfg.RabbitAddr)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer rconn.Close()

	// handle rmq requests
	srv := rmq.NewServer(rconn)
	server.RegisterMathServer(srv, &server.MathImpl{})

	// start the rmq server and schedule a stop
	srv.ServeAsync()
	defer srv.Stop()

	// wait for a program exit to stop the rmq server
	xos.WaitForExit()
}
