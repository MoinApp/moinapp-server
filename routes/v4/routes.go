package v4

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(router *mux.Router) {
	router.StrictSlash(true)

	router.Handle("/moin", defaultHandlerF(serveMoin)).Methods("POST")

	router.Handle("/users/signup", defaultUnauthorizedHandler(http.HandlerFunc(serveSignUp))).Methods("POST")
	router.Handle("/users/auth", defaultUnauthorizedHandler(http.HandlerFunc(serveAuthentication))).Methods("POST")
	router.Handle("/users", defaultHandlerF(serveSearchUser)).Methods("GET")
	router.Handle("/users/recents", defaultHandlerF(serveRecentUsers)).Methods("GET")
	router.Handle("/users/addPush", defaultHandlerF(serveAddPushToken)).Methods("POST")

	router.Handle("/user/{username}", defaultHandlerF(serveGetUserProfile)).Methods("GET")
}

func SetHttpsOnly(httpsOnly bool) {
	setHttpsCheckState(httpsOnly)
}
