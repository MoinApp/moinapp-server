package models

import (
	"github.com/MoinApp/moinapp-server"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB

func DB() *gorm.DB {
	if db == nil {
		//TODO: add production database connection
		dbConnection, err := gorm.Open("sqlite3", "./db.sqlite3")

		if err != nil {
			//TODO: bad coding style
			panic(err)
		}

		db = &dbConnection

		// add database table structs here
		db.AutoMigrate(&User{})
	}

	return db
}

func TestDB() bool {
	return (DB().DB().Ping() == nil)
}
