package main

import (
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/xos"
	"github.com/sknv/microrabbit/app/services/math/cfg"
)

func main() {
	cfg := cfg.Parse()

	conn, err := amqp.Dial(cfg.RabbitURL)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	xos.WaitForExit()
}
