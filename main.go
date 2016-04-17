// The moinapp-server package is the executable for the MoinApp server environment.
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/push"
	"github.com/MoinApp/moinapp-server/routes"
)

const (
	APP_NAME = "MoinApp-Server"
	// TODO: let this be written at compile-time
	APP_VERSION = "feature/go-rewrite"
)

func isProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}

func main() {
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
