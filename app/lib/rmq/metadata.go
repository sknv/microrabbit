package rmq

import (
	"github.com/streadway/amqp"
)

const (
	responseCodeKey = "rmq.responseCode"

	responseOK    int16 = 0
	responseError int16 = 1
)

func HasError(message *amqp.Delivery) bool {
	headers := message.Headers
	code, exist := headers[responseCodeKey]
	if !exist { // if there is no such header, we are ok
		return false
	}
	if code != responseError {
		return false
	}
	return true
}

func WithError(message *amqp.Publishing) *amqp.Publishing {
	if message.Headers == nil {
		message.Headers = make(amqp.Table)
	}
	message.Headers[responseCodeKey] = responseError
	return message
}
