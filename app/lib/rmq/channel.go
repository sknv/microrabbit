package rmq

import (
	"github.com/streadway/amqp"
)

type Channel struct {
	*amqp.Channel
}

func NewChannel(conn *amqp.Connection) (*Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Channel{Channel: ch}, nil
}

func (c *Channel) DeclareQueue(name string, durable bool) (amqp.Queue, error) {
	return c.QueueDeclare(
		name,    // name
		durable, // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
}

func (c *Channel) QoS(prefetchCount int) error {
	return c.Qos(
		prefetchCount, // prefetch count
		0,             // prefetch size
		false,         // global
	)
}

func (c *Channel) ConsumeFrom(queueName string, autoAck bool) (<-chan amqp.Delivery, error) {
	return c.Channel.Consume(
		queueName, // queue
		"",        // consumer
		autoAck,   // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}

func (c *Channel) PublishTo(exchange, routingKey string, message *amqp.Publishing) error {
	return c.Channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		*message,   // message
	)
}
