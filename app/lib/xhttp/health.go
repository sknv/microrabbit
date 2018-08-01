package xhttp

import (
	"net/http"
)

type HealthServer struct{}

func (*HealthServer) Check(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}
