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
	portNum, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint(portNum)
}
