package xamqp

import (
	"fmt"
	"net/http"
)

type ErrorCode uint32

const (
	OK               ErrorCode = 0
	InvalidArgument  ErrorCode = 1
	Unauthenticated  ErrorCode = 2
	PermissionDenied ErrorCode = 3
	Internal         ErrorCode = 4
)

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func NewError(code ErrorCode, message string) *Error {
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

func ServerHTTPStatusFromErrorCode(code ErrorCode) int {
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

func IsValidErrorCode(code ErrorCode) bool {
	return ServerHTTPStatusFromErrorCode(code) != 0
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

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
	return fmt.Sprintf("amqp error %s: %s", e.Code, e.Message)
}
