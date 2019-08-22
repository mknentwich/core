package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Initializes the database with the whole table structure
func InitializeDb() error {
	database, err := sql.Open("sqlite3", "./mknentwich.db")
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
