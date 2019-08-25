package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mknentwich/core/context"
)

var log context.Log

// Initializes the database with the whole table structure
func initializeDb() error {
	database, err := sql.Open("sqlite3", context.Conf.SQLiteFile)
	if err != nil {
		return err
	}
	defer database.Close()

	for _, statement := range createStatements {
		executable, err := database.Prepare(statement)
		if err != nil {
			return err
		}
		executable.Exec()
	}
	return err
}

//Serve call for the service registry
func Serve(logger context.Log) (context.ServiceResult, error) {
	log = logger
	return context.ServiceResult{}, initializeDb()
}
