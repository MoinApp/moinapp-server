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

	router.HandleFunc("/", serve_Root).Methods("GET", "POST")

	router.HandleFunc("/users/signup", serve_Users_SignUp).Methods("POST")
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

func http_redirect(w http.ResponseWriter, newLocation string) {
	w.Header().Add("Location", newLocation)
	w.WriteHeader(http.StatusFound)
}

func serve_Root(w http.ResponseWriter, r *http.Request) {
	http_redirect(w, homeRedirectURL)
}
