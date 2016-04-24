// The moinapp-server package is the executable for the MoinApp server environment.
package main

import (
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/push"
	"github.com/MoinApp/moinapp-server/routes"
)

const (
	APP_NAME = "MoinApp-Server"
)

var (
	ErrIncorrectCompilation = errors.New("The application was not compiled correctly. Please refer to the README for further details.")
)

var APP_VERSION string

func isProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}
func checkCorrectCompilation() {
	if APP_VERSION == "" {
		panic(ErrIncorrectCompilation)
	}
}

func main() {
	checkCorrectCompilation()

	log.Printf("%v %q on %v/%v\n", APP_NAME, APP_VERSION, runtime.GOOS, runtime.GOARCH)
	log.Println("Hello! Booting...")
	models.InitDB(isProduction())
	if !models.TestDB() {
		log.Fatal("Database is not connected.")
	}
	push.InitPushServices(isProduction())
	router := routes.CreateRouter(isProduction())

	listeningDone := make(chan error)
	listeningAddr := routes.StartListening(router, listeningDone)
	log.Printf("Ready. Listening on %q.", listeningAddr)
	log.Fatalf("Error: %v.", <-listeningDone)
}
