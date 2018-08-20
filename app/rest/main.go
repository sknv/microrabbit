package main

import (
	"time"

	"github.com/go-chi/chi"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/xchi"
	"github.com/sknv/microrabbit/app/lib/xhttp"
	"github.com/sknv/microrabbit/app/lib/xos"
	"github.com/sknv/microrabbit/app/rest/cfg"
	"github.com/sknv/microrabbit/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second

	healthCheckURL = "/healthz"
)

func main() {
	cfg := cfg.Parse()

	// connect to RabbitMQ
	rmqConn, err := rmq.DialWithReconnect(cfg.RabbitAddr)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer rmqConn.Close()

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	server.RegisterRestServer(rmqConn, router) // handle requests
	registerHealthServer(router)               // handle health check requests

	// start the http server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(serverShutdownTimeout)

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}

func registerHealthServer(router chi.Router) {
	var health xhttp.HealthServer
	router.Get(healthCheckURL, health.Check)
}
