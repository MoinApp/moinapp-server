package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var db *gorm.DB

func InitDB(isProduction bool) *gorm.DB {
	if db == nil {
		var dbConnection gorm.DB
		var err error

		if isProduction {
			dbConnection, err = gorm.Open("postgres", getDatabaseURL())
		} else {
			dbConnection, err = gorm.Open("sqlite3", "./db.sqlite3")
		}

		if err != nil {
			//TODO: bad coding style
			log.Fatalf("Could not connect to database: %v!\n", err)
		}

		db = &dbConnection

		// add database table structs here
		db.AutoMigrate(&User{})

		db.LogMode(!isProduction)
	}

	return db
}

func getDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func TestDB() bool {
	return (DB().DB().Ping() == nil)
}
