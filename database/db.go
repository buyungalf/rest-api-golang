package database

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// github.com/denisenkom/go-mssqldb

var Db *gorm.DB

func InitDb() *gorm.DB { // OOP constructor
	Db = connectDB()
	return Db
}

func connectDB() (*gorm.DB) {	
	dsn := "sqlserver://buyung123:buyung123@localhost:1433?database=shop"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err !=nil {
		fmt.Println("Error...")
		return nil
	}
	return db
}