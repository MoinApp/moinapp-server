package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var db *gorm.DB

// InitDB initializes the database connection and returns the handle. This may crash if the connection cannot be established.
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
			return nil
		}

		db = &dbConnection

		// add database table structs here
		db.AutoMigrate(&User{}, &PushToken{})

		db.LogMode(!isProduction)
	}

	return db
}

func getDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

// TestDB checks if the database connected is responding to a PING.
func TestDB() bool {
	if db == nil {
		return false
	}
	return (db.DB().Ping() == nil)
}
