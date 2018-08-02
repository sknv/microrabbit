package xhttp

import (
	"net/http"
)

func AbortHandler() {
	panic(http.ErrAbortHandler)
}

func AbortHandlerWithInternalError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	AbortHandler()
}

func IsHandlerAborted(v interface{}) bool {
	return v == http.ErrAbortHandler
}
