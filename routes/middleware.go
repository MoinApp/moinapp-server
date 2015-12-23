package routes

// BUG(sgade): The middleware errors are plain text. They should be JSON APIErrors.

import (
	"compress/gzip"
	"io"
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

type gzipResponseWriter struct {
	// Writer is the gzip compressing io.Writer.
	io.Writer
	// ResponseWriter is the standard response writer for the http connection.
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

// setHttpsCheckState sets the flag whether any request should be checked to be made via the HTTPS-protocol
func setHttpsCheckState(httpsCheckEnabled bool) {
	httpsOnlyCheckEnabled = httpsCheckEnabled
}

func defaultHandlerF(nextFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	next := http.HandlerFunc(nextFunc)
	return defaultHandler(next)
}

func defaultHandler(next http.Handler) http.Handler {
	return httpsCheckHandler(gzipCompressionHandler(securityHandler(timeoutHandler(headerHandler(next)))))
}

func defaultUnauthorizedHandler(next http.Handler) http.Handler {
	return httpsCheckHandler(gzipCompressionHandler(timeoutHandler(headerHandler(next))))
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

func gzipCompressionHandler(next http.Handler) http.Handler {
	// with help from https://gist.github.com/the42/1956518
	fn := func(rw http.ResponseWriter, req *http.Request) {
		if !strings.Contains(strings.ToLower(req.Header.Get("Accept-Encoding")), "gzip") {
			next.ServeHTTP(rw, req)
			return
		}

		rw.Header().Set("Content-Encoding", "gzip")

		compressor := gzip.NewWriter(rw)
		defer compressor.Close()
		newWriter := gzipResponseWriter{
			Writer:         compressor,
			ResponseWriter: rw,
		}

		next.ServeHTTP(newWriter, req)
	}

	return http.HandlerFunc(fn)
}

func timeoutHandler(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, timeout, "Response timeout reached.")
}

func headerHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.Header().Set("X-Served-by", "moinapp-server")

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