package main

import (
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/xos"
	"github.com/sknv/microrabbit/app/math/cfg"
)

func main() {
	cfg := cfg.Parse()

	// connect to RabbitMQ
	rconn, err := amqp.Dial(cfg.RabbitAddr)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer rconn.Close()

	// // start the rmq server and schedule a stop
	// srv := server.NewMathServer(rconn)
	// xos.FailOnError(srv.ServeAsync(), "failed to start a math server")
	// defer srv.Stop()

	// // handle nats requests
	// natsServer := xnats.NewServer(natsConn)
	// server.RegisterMathServer(natsServer, &server.MathImpl{})

	// log.Print("[INFO] math service started")
	// defer log.Print("[INFO] math service stopped")

	xos.WaitForExit()
}
