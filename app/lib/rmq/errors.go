package rmq

import (
	"fmt"
	"net/http"

	"github.com/streadway/amqp"
)

type StatusCode uint32

const (
	StatusOK               StatusCode = 0
	StatusInvalidArgument  StatusCode = 1
	StatusUnauthenticated  StatusCode = 2
	StatusPermissionDenied StatusCode = 3
	StatusInternal         StatusCode = 4
	StatusDeadlineExceeded StatusCode = 5
)

func NewError(code StatusCode, message string) *Error {
	if IsValidErrorCode(code) {
		return &Error{
			Code:    uint32(code),
			Message: message,
		}
	}
	return &Error{
		Code:    uint32(StatusInternal),
		Message: "invalid error code: " + fmt.Sprint(code),
	}
}

func WrapError(code StatusCode, err error) *Error {
	return NewError(code, err.Error())
}

func FromError(err error) (*Error, bool) {
	qerr, ok := err.(*Error)
	return qerr, ok
}

func ServerHTTPStatusFromErrorCode(code StatusCode) int {
	switch code {
	case StatusOK:
		return http.StatusOK
	case StatusInvalidArgument:
		return http.StatusBadRequest
	case StatusUnauthenticated:
		return http.StatusUnauthorized
	case StatusPermissionDenied:
		return http.StatusForbidden
	case StatusInternal:
		return http.StatusInternalServerError
	case StatusDeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return 0 // invalid
	}
}

func IsValidErrorCode(code StatusCode) bool {
	return ServerHTTPStatusFromErrorCode(code) != 0
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

type headerCode uint16

const (
	headerCodeKey = "code"

	headerOK    headerCode = 0
	headerError headerCode = 1
)

func MessageHasError(message *amqp.Delivery) bool {
	headers := message.Headers
	code, ok := headers[headerCodeKey]
	if !ok { // if there is no such header, we are ok
		return false
	}
	if code != headerError {
		return false
	}
	return true
}

func MessageWithError(message *amqp.Publishing) *amqp.Publishing {
	message.Headers[headerCodeKey] = headerError
	return message
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func (e *Error) StatusCode() StatusCode {
	if e == nil {
		return StatusOK
	}
	return StatusCode(e.Code)
}

func (e *Error) MetaValue(key string) string {
	if e.Meta != nil {
		return e.Meta[key] // also returns "" if key is not in meta map
	}
	return ""
}

func (e *Error) WithMeta(key string, value string) *Error {
	newErr := &Error{
		Code:    e.Code,
		Message: e.Message,
		Meta:    make(map[string]string, len(e.Meta)),
	}
	for key, val := range e.Meta { // copy existing map
		newErr.Meta[key] = val
	}
	newErr.Meta[key] = value // upsert the value
	return newErr
}

func (e *Error) Error() string {
	return fmt.Sprintf("amqp error %d: %s", e.Code, e.Message)
}
