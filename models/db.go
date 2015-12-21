package models

import (
	"github.com/MoinApp/moinapp-server/global"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *gorm.DB

func DB() *gorm.DB {
	if db == nil {
		var dbConnection gorm.DB
		var err error

		if global.IsProduction() {
			dbConnection, err = gorm.Open("postgres", global.GetDatabaseURL())
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

		db.LogMode(!global.IsProduction())
	}

	return db
}

func TestDB() bool {
	return (DB().DB().Ping() == nil)
}
