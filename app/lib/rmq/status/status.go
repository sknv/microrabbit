package status

import (
	"fmt"
	"net/http"
)

type statusCode uint32

const (
	OK               statusCode = 0
	InvalidArgument  statusCode = 1
	Unauthenticated  statusCode = 2
	PermissionDenied statusCode = 3
	Internal         statusCode = 4
	DeadlineExceeded statusCode = 5
)

func ServerHTTPStatusFromErrorCode(code statusCode) int {
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
	case DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return 0 // invalid
	}
}

func IsValidErrorCode(code statusCode) bool {
	return ServerHTTPStatusFromErrorCode(code) != 0
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func Error(code statusCode, message string) *Status {
	if IsValidErrorCode(code) {
		return &Status{
			Code:    uint32(code),
			Message: message,
		}
	}
	return &Status{
		Code:    uint32(Internal),
		Message: "rmq: invalid status code: " + fmt.Sprint(code),
	}
}

func FromError(err error) (*Status, bool) {
	stat, match := err.(*Status)
	if match {
		return stat, true
	}
	return Error(Internal, err.Error()), false
}

func (s *Status) StatusCode() statusCode {
	if s == nil {
		return OK
	}
	return statusCode(s.Code)
}

func (s *Status) MetaValue(key string) string {
	if s.Meta != nil {
		return s.Meta[key] // also returns "" if key is not in meta map
	}
	return ""
}

func (s *Status) WithMeta(key string, value string) *Status {
	newErr := &Status{
		Code:    s.Code,
		Message: s.Message,
		Meta:    make(map[string]string, len(s.Meta)),
	}
	for key, val := range s.Meta { // copy existing map
		newErr.Meta[key] = val
	}
	newErr.Meta[key] = value // upsert the value
	return newErr
}

func (s *Status) Error() string {
	return fmt.Sprintf("rmq error %d: %s", s.Code, s.Message)
}
