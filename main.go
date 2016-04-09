// The moinapp-server package is the executable for the MoinApp server environment.
package main

import (
	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/push"
	"github.com/MoinApp/moinapp-server/routes"
	"log"
	"os"
)

func isProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}

func main() {
	log.Println("Hello! Booting...")
	models.InitDB(isProduction())
	if !models.TestDB() {
		log.Fatal("Database is not connected.")
	}
	push.InitPushServices(isProduction())
	router := routes.CreateRouter(isProduction())

	listeningDone := make(chan bool)
	listeningAddr := routes.StartListening(router, listeningDone)
	log.Printf("Ready. Listening on %q.", listeningAddr)
	<-listeningDone
}
