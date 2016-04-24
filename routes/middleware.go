package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
)

const (
	// request response timeout
	timeout = 5000 * time.Millisecond
)

func middleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, middleware_timeout(handlers.CompressHandler(middleware_defaultHeaders(next))))
}

// --- --- --- Header --- --- ---

func middleware_defaultHeaders(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("X-Served-by", "moinapp-server")

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// --- --- --- Timeout --- --- ---

func middleware_timeout(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, timeout, "Response timeout reached.")
}
