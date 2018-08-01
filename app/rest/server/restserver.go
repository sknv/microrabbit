package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/streadway/amqp"

	"github.com/sknv/microrabbit/app/lib/xamqp"
	"github.com/sknv/microrabbit/app/lib/xhttp"
	math "github.com/sknv/microrabbit/app/services/math/public"
)

type RestServer struct {
	mathClient *math.MathClient
}

func NewRestServer(conn *amqp.Connection) *RestServer {
	return &RestServer{
		mathClient: math.NewClient(conn),
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

	reply, err := s.mathClient.Circle(context.Background(), &args)
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

	// check if the error is an *xamqp.Error
	qerr, ok := xamqp.FromError(err)
	if !ok {
		log.Print("[ERROR] abort on error: ", err)
		panic(err)
	}

	status := xamqp.ServerHTTPStatusFromErrorCode(qerr.StatusCode())
	if status != http.StatusInternalServerError {
		log.Print("[ERROR] abort on error: ", qerr)
		http.Error(w, qerr.Message, status)
		xhttp.AbortHandler()
	}
	panic(qerr)
}
