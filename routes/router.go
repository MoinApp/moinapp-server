package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MoinApp/moinapp-server/routes/v4"
	"github.com/gorilla/mux"
)

const (
	homeRedirectURL = "https://i.imgur.com/E2T98iu.jpg"
	defaultPort     = 3000
)

func CreateRouter(httpsOnly bool) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/", http.RedirectHandler(homeRedirectURL, http.StatusFound)).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.PathPrefix("/v1").HandlerFunc(discontinuationHandler)
	apiRouter.PathPrefix("/v2").HandlerFunc(discontinuationHandler)
	apiRouter.PathPrefix("/v3").HandlerFunc(discontinuationHandler)
	v4.SetHttpsOnly(httpsOnly)
	v4.RegisterRoutes(apiRouter.PathPrefix("/v4").Subrouter())
	// fallback, for old api calls (pre-v4)
	apiRouter.HandleFunc("/", discontinuationHandler)

	return router
}

func getListeningPort() uint {
	port := os.Getenv("PORT")
	portNum, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		if len(port) > 0 {
			defaultPortName := fmt.Sprintf("%v", defaultPort)
			if defaultPort == 0 {
				defaultPortName = "system defined port"
			}

			log.Printf("Error parsing PORT. Using %v.", defaultPortName)
		}

		portNum = defaultPort
	}
	return uint(portNum)
}

func StartListening(router *mux.Router, listeningError chan error) string {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%v", getListeningPort()),
		Handler: middleware(router),
	}

	go func() {
		listeningError <- srv.ListenAndServe()
	}()

	return srv.Addr
}
