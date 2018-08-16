package xhttp

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return &Server{Server: srv}
}

func (s *Server) ListenAndServeAsync() {
	log.Print("[INFO] starting the http server on ", s.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Print("[ERROR] failed to serve the http server: ", err)
		}
	}()
}

func (s *Server) StopGracefully(shutdownTimeout time.Duration) {
	log.Print("[INFO] stopping the http server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("[FATAL] failed to stop the http server gracefully: ", err)
	}
	log.Print("[INFO] http server gracefully stopped")
}
