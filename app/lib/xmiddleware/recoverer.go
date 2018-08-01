package xmiddleware

import (
	"net/http"

	"github.com/sknv/microrabbit/app/lib/xhttp"
)

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if xhttp.IsHandlerAborted(rvr) {
					return // response is already flushed
				}
				panic(rvr) // throw the error
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
