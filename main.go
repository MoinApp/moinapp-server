// The moinapp-server package is the executable for the MoinApp server environment.
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/MoinApp/moinapp-server/info"
	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/push"
	"github.com/MoinApp/moinapp-server/routes"
)

func isProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}

func main() {
	info.CheckCorrectCompilation()

	log.Printf("%v %q on %v/%v\n", info.AppName, info.AppVersion, runtime.GOOS, runtime.GOARCH)
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
