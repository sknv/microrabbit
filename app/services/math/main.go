package main

import (
	"log"

	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/xos"
	"github.com/sknv/microrabbit/app/services/math/cfg"
)

func main() {
	cfg := cfg.Parse()

	rconn, err := amqp.Dial(cfg.RabbitAddr)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer rconn.Close()

	log.Print("[INFO] starting a rabbit worker")
	xos.WaitForExit()
}
