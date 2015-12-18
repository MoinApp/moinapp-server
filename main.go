package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home)

	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}

func http_redirect(w http.ResponseWriter, newLocation string) {
	w.Header().Add("Location", newLocation)
	w.WriteHeader(http.StatusFound)
}

func home(w http.ResponseWriter, r *http.Request) {
	http_redirect(w, "http://i.imgur.com/E2T98iu.jpg")
}

func users_signUp(w http.ResponseWriter, r *http.Request) {

}
