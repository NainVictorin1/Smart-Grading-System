package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var database *sql.DB

func initdatabase() {
	var err error
	connStr := "postgres://feedback:tapirhorse@localhost/tapir_nain?sslmode=disable"
	database, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("fail toconnect to database", err)
		if err := database.Ping(); err != nil {
			log.Fatal("failed to ping the database", err)
		}
		log.Println("Successfully connected to the database.")
	}
}
