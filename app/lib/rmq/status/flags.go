package status

import (
	"github.com/streadway/amqp"
)

type headerCode uint16

const (
	headerCodeKey = "code"

	headerOK    headerCode = 0
	headerError headerCode = 1
)

func HasError(message *amqp.Delivery) bool {
	headers := message.Headers
	code, exist := headers[headerCodeKey]
	if !exist { // if there is no such header, we are ok
		return false
	}
	if code != headerError {
		return false
	}
	return true
}

func WithError(message *amqp.Publishing) *amqp.Publishing {
	message.Headers[headerCodeKey] = headerError
	return message
}
