package xchi

import (
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/sknv/microrabbit/app/lib/xmiddleware"
)

func UseDefaultMiddleware(router chi.Router) {
	router.Use(
		middleware.RealIP, middleware.Logger, middleware.Recoverer, xmiddleware.Recoverer,
	)
}

func UseThrottle(router chi.Router, limit int) {
	router.Use(middleware.Throttle(limit))
}

func WithThrottle(router chi.Router, limit int) chi.Router {
	return router.With(middleware.Throttle(limit))
}

func UseTimeout(router chi.Router, timeout time.Duration) {
	router.Use(middleware.Timeout(timeout))
}

func WithTimeout(router chi.Router, timeout time.Duration) chi.Router {
	return router.With(middleware.Timeout(timeout))
}

func UseLimitHandler(router chi.Router, limit float64) {
	limiter := tollbooth.NewLimiter(limit, nil)
	router.Use(tollbooth_chi.LimitHandler(limiter))
}

func WithLimitHandler(router chi.Router, limit float64) chi.Router {
	limiter := tollbooth.NewLimiter(limit, nil)
	return router.With(tollbooth_chi.LimitHandler(limiter))
}
