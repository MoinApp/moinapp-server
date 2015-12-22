// The moinapp-server package is the executable for the MoinApp server environment.
package main

import (
	"github.com/MoinApp/moinapp-server/models"
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
	routes.InitRouter()

	log.Println("Ready.")
	routes.StartListening()
}
