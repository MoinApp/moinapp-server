// The moinapp-server package is the executable for the MoinApp server environment.
package main

import (
	"fmt"
	"github.com/MoinApp/moinapp-server/models"
	"github.com/MoinApp/moinapp-server/routes"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Printf("Starting...")
	models.TestDB()
	routes.InitRouter()

	fmt.Printf("\nReady.\n")
	routes.StartListening()
}
