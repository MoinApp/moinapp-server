package global

import "os"

func IsProduction() bool {
	return (os.Getenv("PRODUCTION") != "")
}

func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}
