package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mknentwich/core/context"
)

//log function for the database package
var log context.Log

//database instance pointer
var database *gorm.DB

// Initializes the database with the whole table structure
func initializeDb() error {
	db, err := gorm.Open("sqlite3", context.Conf.SQLiteFile)
	if err != nil {
		return err
	}
	database = db.Exec("pragma foreign_keys = on;").AutoMigrate(&Address{}, &Category{}, &Order{}, &Score{}, &User{}).Set("gorm:auto_preload", true)
	return err
}

//Receive the database instance pointer
func Receive() *gorm.DB {
	return database
}

//Serve call for the service registry
func Serve(logger context.Log) (context.ServiceResult, error) {
	log = logger
	return context.ServiceResult{}, initializeDb()
}
