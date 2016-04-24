package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MoinApp/moinapp-server/info"
	"github.com/gorilla/handlers"
)

const (
	// request response timeout
	timeout = 5000 * time.Millisecond
)

func middleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(middleware_timeout(middleware_recovery(middleware_defaultHeaders(next)))))
}

// --- --- --- Header --- --- ---

func middleware_defaultHeaders(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		servedByHeader := fmt.Sprintf("%v (%v)", info.AppName, info.AppVersion)
		rw.Header().Set("X-Served-by", servedByHeader)

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

// --- --- --- Timeout --- --- ---

func middleware_timeout(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, timeout, "Response timeout reached.")
}

// --- --- --- Panic Recovery --- --- ---

func middleware_recovery(next http.Handler) http.Handler {
	return handlers.RecoveryHandler()(next)
}
