package main

import (
	"os"
	"strconv"
)

func IsProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}

func GetListeningPort() int {
	port := os.Getenv("PORT")
	portNum, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	return portNum
}
