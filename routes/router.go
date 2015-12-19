package routes

import (
	"github.com/gorilla/mux"
	"net/http"
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

func StartListening() {
	http.Handle("/", router)
	http.ListenAndServe(":3000", nil)
}

func http_redirect(w http.ResponseWriter, newLocation string) {
	w.Header().Add("Location", newLocation)
	w.WriteHeader(http.StatusFound)
}

func serve_Root(w http.ResponseWriter, r *http.Request) {
	http_redirect(w, homeRedirectURL)
}
