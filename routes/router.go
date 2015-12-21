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

func InitRouter() {
	router = mux.NewRouter()

	router.HandleFunc("/", serveRoot).Methods("GET", "POST")

	router.HandleFunc("/moin", serveMoin).Methods("POST")

	router.HandleFunc("/users/signup", serveSignUp).Methods("POST")
	router.HandleFunc("/users/auth", serveAuthentication).Methods("POST")
	router.HandleFunc("/users", serveSearchUser).Methods("GET")
	router.HandleFunc("/users/recents", serveRecentUsers).Methods("GET")
	router.HandleFunc("/users/addPush", serveAddPushToken).Methods("POST")

	router.HandleFunc("/user/{username}", serveGetUserProfile).Methods("GET")
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

func httpRedirect(w http.ResponseWriter, newLocation string) {
	w.Header().Add("Location", newLocation)
	w.WriteHeader(http.StatusFound)
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	httpRedirect(w, homeRedirectURL)
}
