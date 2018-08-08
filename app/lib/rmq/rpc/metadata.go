package rpc

import (
	"github.com/streadway/amqp"
)

type responseCode uint16

const (
	responseCodeKey = "rmq.responseCode"

	responseOK    responseCode = 0
	responseError responseCode = 1
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
	message.Headers[responseCodeKey] = responseError
	return message
}
