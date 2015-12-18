package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
)

var (
	isProduction = false
	db           gorm.DB
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	isProduction = (os.Getenv("PRODUCTION") != "")

	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/", home)

	http.Handle("/", r)
	fmt.Printf("Serving...\n")
	http.ListenAndServe(":3000", nil)
}

func initDB() {
	var err error

	if isProduction {
		// heroku config
		dbURL := ""
		db, err = gorm.Open("postgres", dbURL)
	} else {
		db, err = gorm.Open("sqlite3", "./db.sqlite3")
	}
	if err != nil {
		handleError(err)
	}

	err = db.DB().Ping()
	if err != nil {
		fmt.Printf("Db error: %v\n", err)
	} else {
		fmt.Printf("Database ready.\n")
	}
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
