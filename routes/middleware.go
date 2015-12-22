package routes

// BUG(sgade): The middleware errors are plain text. They should be JSON APIErrors.

import (
	"net/http"
	"strings"
	"time"
)

const (
	timeout = 1000 * time.Millisecond
)

var (
	httpsOnlyCheckEnabled = true
)

// setHttpsCheckState sets the flag whether any request should be checked to be made via the HTTPS-protocol
func setHttpsCheckState(httpsCheckEnabled bool) {
	httpsOnlyCheckEnabled = httpsCheckEnabled
}

func defaultHandlerF(nextFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	next := http.HandlerFunc(nextFunc)
	return defaultHandler(next)
}

func defaultHandler(next http.Handler) http.Handler {
	return httpsCheckHandler(securityHandler(defaultTimeoutHandler(defaultHeaderHandler(next))))
}

func defaultUnauthorizedHanldler(next http.Handler) http.Handler {
	return httpsCheckHandler(defaultTimeoutHandler(defaultHeaderHandler(next)))
}

func httpsCheckHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		// only check if it is enabled
		if httpsOnlyCheckEnabled {
			// get the index in "HTTP/1.1"
			slashIndex := strings.Index(req.Proto, "/")

			// maybe we got something different than a HTTP protocol?
			if slashIndex != -1 {
				protocol := strings.ToLower(req.Proto[:slashIndex])

				if protocol != "https" {
					data := []byte("Only https allowed")
					rw.Write(data)
					return
				}
			}
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
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

func securityHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Session")

		// TODO: Check
		if token == "" {
			data := []byte("Authentication required.")
			rw.Write(data)
			return
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}
