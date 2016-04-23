package v4

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/routes/v4/auth"
)

const (
	requestUserHeader = "_moinapp_user"
)

var (
	httpsOnlyCheckEnabled = true
)

var (
	ErrOnlyHttpsAllowed = errors.New("Only https allowed")
)

// getUserFromRequest returns the user model given a valid request header.
func getUserFromRequest(req *http.Request) *models.User {
	userID := req.Header.Get(requestUserHeader)

	if len(userID) == 0 {
		return nil
	}

	return models.FindUserById(userID)
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
	return httpsCheckHandler(securityHandler(headerHandler(next)))
}

func defaultUnauthorizedHandler(next http.Handler) http.Handler {
	return httpsCheckHandler(headerHandler(next))
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
					sendErrorCode(rw, ErrOnlyHttpsAllowed, http.StatusForbidden)
					return
				}
			}
		}

		next.ServeHTTP(rw, req)
	}

	return http.HandlerFunc(fn)
}

func headerHandler(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

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
