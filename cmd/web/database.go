package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Database *sql.DB

// initDatabase initializes the connection to the database
func initDatabase(dsn string) error {
	var err error
	// Open the database connection using the provided DSN
	Database, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	// Ensure the connection is successful by pinging the database
	if err := Database.Ping(); err != nil {
		return err
	}

	log.Println("Successfully connected to the database.")
	return nil
}
