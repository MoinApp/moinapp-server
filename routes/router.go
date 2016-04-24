package routes

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MoinApp/moinapp-server/routes/v4"
	"github.com/gorilla/mux"
)

const (
	homeRedirectURL        = "https://i.imgur.com/E2T98iu.jpg"
	defaultPort     uint16 = 3000
	listenNet              = "tcp"
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
	apiRouter.NotFoundHandler = http.HandlerFunc(discontinuationHandler)

	return router
}

func getListeningPort() uint16 {
	port := strings.Trim(os.Getenv("PORT"), " ")
	if len(port) > 0 {
		if portNum, err := strconv.ParseUint(port, 10, 16); err == nil {
			return uint16(portNum)
		} else {
			defaultPortName := fmt.Sprintf("%v", defaultPort)
			if defaultPort == 0 {
				defaultPortName = "system defined port"
			}

			log.Printf("Error parsing PORT. Using %v.", defaultPortName)
		}
	}

	return defaultPort
}

func StartListening(router *mux.Router, listeningError chan error) net.Addr {
	srv := http.Server{
		Handler: middleware(router),
	}
	addr, err := net.ResolveTCPAddr(listenNet, fmt.Sprintf(":%v", getListeningPort()))
	if err != nil {
		log.Fatalf("Could not parse %v addr: %v.", listenNet, err)
	}
	listener, err := net.ListenTCP(listenNet, addr)
	if err != nil {
		log.Fatalf("Could not listen on %v: %v.", addr, err)
	}

	go func() {
		listeningError <- srv.Serve(listener)
	}()

	return listener.Addr()
}
