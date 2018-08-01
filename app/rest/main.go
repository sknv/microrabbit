package main

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/xchi"
	"github.com/sknv/microrabbit/app/lib/xhttp"
	"github.com/sknv/microrabbit/app/lib/xos"
	"github.com/sknv/microrabbit/app/rest/cfg"
	"github.com/sknv/microrabbit/app/rest/server"
)

const (
	concurrentRequestLimit = 1000
	serverShutdownTimeout  = 60 * time.Second
)

func main() {
	cfg := cfg.Parse()

	conn, err := amqp.Dial(cfg.RabbitURL)
	xos.FailOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	// config the http router
	router := chi.NewRouter()
	xchi.UseDefaultMiddleware(router)
	xchi.UseThrottle(router, concurrentRequestLimit)

	// handle requests
	rest := server.NewRestServer(conn)
	rest.Route(router)

	// handle health check requests
	var health xhttp.HealthServer
	router.Get("/healthz", health.Check)

	// start the http server and schedule a stop
	srv := xhttp.NewServer(cfg.Addr, router)
	srv.ListenAndServeAsync()
	defer srv.StopGracefully(serverShutdownTimeout)

	// wait for a program exit to stop the http server
	xos.WaitForExit()
}
