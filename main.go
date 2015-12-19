package main

import (
	"fmt"
	"github.com/MoinApp/moinapp-server/models"
	"github.com/gorilla/mux"

	"net/http"
	"os"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	models.TestDB()

	r := mux.NewRouter()
	r.HandleFunc("/", home)

	http.Handle("/", r)
	fmt.Printf("Serving...\n")
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
