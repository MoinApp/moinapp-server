// The moinapp-server package is the executable for the MoinApp server environment.
package main

import (
	"fmt"
	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/routes"
	"os"
)

func isProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}

func main() {
	fmt.Printf("Starting...")
	models.TestDB()
	models.InitDB(isProduction())
	routes.InitRouter()

	fmt.Printf("\nReady.\n")
	routes.StartListening()
}
