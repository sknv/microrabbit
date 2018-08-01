package xhttp

import (
	"net/http"
)

func AbortHandler() {
	panic(http.ErrAbortHandler)
}

func IsHandlerAborted(v interface{}) bool {
	return v == http.ErrAbortHandler
}
