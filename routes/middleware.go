package routes

import (
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
)

const (
	// request response timeout
	timeout = 5000 * time.Millisecond
)

func middleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, middleware_timeout(middleware_gzip(middleware_defaultHeaders(next))))
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

// --- --- --- GZIP Compression --- --- ---

// gzipResponseWriter handles GZIP responses
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

func middleware_gzip(next http.Handler) http.Handler {
	// with help from https://gist.github.com/the42/1956518
	fn := func(rw http.ResponseWriter, req *http.Request) {
		// only send gzip if supported by the requesting client
		if !strings.Contains(strings.ToLower(req.Header.Get("Accept-Encoding")), "gzip") {
			// if its not, then serve normal request
			next.ServeHTTP(rw, req)
			return
		}

		// add content-encoding header
		rw.Header().Set("Content-Encoding", "gzip")

		// create compressor for this request
		compressor := gzip.NewWriter(rw)
		// close this gzip writer after ServeHTTP returns
		defer compressor.Close()
		newWriter := gzipResponseWriter{
			Writer:         compressor,
			ResponseWriter: rw,
		}

		// serve with new gzip compressor
		next.ServeHTTP(newWriter, req)
	}

	return http.HandlerFunc(fn)
}
