package rmq

import (
	"fmt"
	"net/http"
)

type StatusCode uint32

const (
	OK               StatusCode = 0
	InvalidArgument  StatusCode = 1
	Unauthenticated  StatusCode = 2
	PermissionDenied StatusCode = 3
	Internal         StatusCode = 4
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func NewError(code StatusCode, message string) *Error {
	if IsValidErrorCode(code) {
		return &Error{
			Code:    uint32(code),
			Message: message,
		}
	}
	return &Error{
		Code:    uint32(Internal),
		Message: "invalid error code: " + fmt.Sprint(code),
	}
}

func FromError(err error) (*Error, bool) {
	qerr, ok := err.(*Error)
	return qerr, ok
}

func ServerHTTPStatusFromErrorCode(code StatusCode) int {
	switch code {
	case OK:
		return http.StatusOK
	case InvalidArgument:
		return http.StatusBadRequest
	case Unauthenticated:
		return http.StatusUnauthorized
	case PermissionDenied:
		return http.StatusForbidden
	case Internal:
		return http.StatusInternalServerError
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

func (e *Error) StatusCode() StatusCode {
	if e == nil {
		return OK
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
