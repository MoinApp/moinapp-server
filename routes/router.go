package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

const (
	homeRedirectURL = "http://i.imgur.com/E2T98iu.jpg"
)

var router *mux.Router

func InitRouter(httpsOnly bool) {
	setHttpsCheckState(httpsOnly)
	router = mux.NewRouter()

	router.Handle("/", defaultHandler(http.RedirectHandler(homeRedirectURL, http.StatusFound))).Methods("GET", "POST")

	router.Handle("/moin", defaultHandlerF(serveMoin)).Methods("POST")

	router.Handle("/users/signup", defaultUnauthorizedHanldler(http.HandlerFunc(serveSignUp))).Methods("POST")
	router.Handle("/users/auth", defaultUnauthorizedHanldler(http.HandlerFunc(serveAuthentication))).Methods("POST")
	router.Handle("/users", defaultHandlerF(serveSearchUser)).Methods("GET")
	router.Handle("/users/recents", defaultHandlerF(serveRecentUsers)).Methods("GET")
	router.Handle("/users/addPush", defaultHandlerF(serveAddPushToken)).Methods("POST")

	router.Handle("/user/{username}", defaultHandlerF(serveGetUserProfile)).Methods("GET")
}

func getListeningPort() uint {
	port := os.Getenv("PORT")
	portNum, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		portNum = 3000
	}
	return uint(portNum)
}

func StartListening() {
	listenFormat := fmt.Sprintf(":%v", getListeningPort())

	http.Handle("/", router)
	http.ListenAndServe(listenFormat, nil)
}
