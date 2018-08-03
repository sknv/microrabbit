package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/rmq"
	"github.com/sknv/microrabbit/app/lib/xhttp"
	math "github.com/sknv/microrabbit/app/services/math/rpc"
)

type RestServer struct {
	mathClient *math.MathClient
}

func NewRestServer(rconn *amqp.Connection) *RestServer {
	return &RestServer{
		mathClient: math.NewClient(rconn),
	}
}

func (s *RestServer) Route(router chi.Router) {
	router.Route("/math", func(r chi.Router) {
		r.Get("/circle", s.Circle)
	})
}

func (s *RestServer) Circle(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	radius := parseFloat(w, queryParams.Get("r"))
	args := math.CircleArgs{
		Radius: radius,
	}

	// set the reply timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reply, err := s.mathClient.Circle(ctx, &args)
	abortOnError(w, err)
	render.JSON(w, r, reply)
}

// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------
// ----------------------------------------------------------------------------

func parseFloat(w http.ResponseWriter, s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Print("[ERROR] parse float: ", err)
		http.Error(w, "argument must be a float number", http.StatusBadRequest)
		xhttp.AbortHandler()
	}
	return val
}

func abortOnError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	log.Print("[ERROR] abort on error: ", err)

	// check if the error is an *rmq.Error
	cause := errors.Cause(err)
	rerr, ok := rmq.FromError(cause)
	if !ok {
		xhttp.AbortHandlerWithInternalError(w)
	}

	status := rmq.ServerHTTPStatusFromErrorCode(rerr.StatusCode())
	if status != http.StatusInternalServerError {
		http.Error(w, rerr.GetMessage(), status)
		xhttp.AbortHandler()
	}
	xhttp.AbortHandlerWithInternalError(w)
}
