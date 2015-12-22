package routes

import (
	"net/http"
	"time"
)

const (
	timeout = 1000 * time.Millisecond
)

func defaultHandlerF(nextFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	next := http.HandlerFunc(nextFunc)
	return defaultHandler(next)
}

func defaultHandler(next http.Handler) http.Handler {
	return defaultTimeoutHandler(defaultHeaderHandler(next))
}

func defaultTimeoutHandler(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, timeout, "Response timeout reached.")
}

func defaultHeaderHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		rw.Header().Add("X-Served-by", "moinapp-server")

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}
