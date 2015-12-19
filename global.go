package main

import "os"

func IsProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}
