package v4

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/routes/v4/auth"
	"github.com/gorilla/handlers"
)

const (
	timeout           = 1000 * time.Millisecond
	requestUserHeader = "_moinapp_user"
)

var (
	httpsOnlyCheckEnabled = true
)

var (
	ErrOnlyHttpAllowed = errors.New("Only http allowed")
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
	return httpsCheckHandler(gzipCompressionHandler(handlers.LoggingHandler(os.Stdout, securityHandler(timeoutHandler(headerHandler(next))))))
}

func defaultUnauthorizedHandler(next http.Handler) http.Handler {
	return httpsCheckHandler(gzipCompressionHandler(handlers.LoggingHandler(os.Stdout, timeoutHandler(headerHandler(next)))))
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
					sendErrorCode(rw, ErrOnlyHttpAllowed, http.StatusForbidden)
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
		user, err := auth.ValidateSession(req)
		if err != nil || !user.IsResult() {
			sendErrorCode(rw, err, http.StatusForbidden)
			return
		}

		req.Header.Add(requestUserHeader, strconv.Itoa(int(user.ID)))

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

func getUserFromRequest(req *http.Request) *models.User {
	userID := req.Header.Get(requestUserHeader)

	if len(userID) == 0 {
		return nil
	}

	return models.FindUserById(userID)
}
