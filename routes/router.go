package routes

import (
	"fmt"
	"github.com/MoinApp/moinapp-server/routes/v4"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

const (
	homeRedirectURL = "https://i.imgur.com/E2T98iu.jpg"
)

func CreateRouter(httpsOnly bool) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/", http.RedirectHandler(homeRedirectURL, http.StatusFound)).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.PathPrefix("/").HandlerFunc(discontinuationHandler)
	apiRouter.PathPrefix("/v1").HandlerFunc(discontinuationHandler)
	apiRouter.PathPrefix("/v2").HandlerFunc(discontinuationHandler)
	apiRouter.PathPrefix("/v3").HandlerFunc(discontinuationHandler)
	v4.SetHttpsOnly(httpsOnly)
	v4.RegisterRoutes(apiRouter.PathPrefix("/v4").Subrouter())

	return router
}

func getListeningPort() uint {
	port := os.Getenv("PORT")
	portNum, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		portNum = 3000
	}
	return uint(portNum)
}

func StartListening(router *mux.Router, listeningError chan error) string {
	listenFormat := fmt.Sprintf(":%v", getListeningPort())

	srv := http.Server{
		Addr:    listenFormat,
		Handler: router,
	}

	go func() {
		listeningError <- srv.ListenAndServe()
	}()

	return listenFormat
}
