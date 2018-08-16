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
	conn, err := rmq.Dial(cfg.RabbitAddr)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	// handle rmq requests
	srv := rmq.NewServer(conn, interceptors.WithLogger)
	server.RegisterMathServer(srv, &server.MathImpl{})

	// start the rmq server and schedule a stop
	srv.ServeAsync()
	defer srv.Stop()

	// wait for a program exit to stop the rmq server
	xos.WaitForExit()
}
