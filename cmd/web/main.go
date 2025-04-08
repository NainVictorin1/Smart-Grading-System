package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the database connection
	initdatabase()
	defer database.Close() // Ensure database closes on exit

	// Register routes
	registerRoutes()

	log.Println("Server started on :4000")
	log.Fatal(http.ListenAndServe(":4000", nil))

}
