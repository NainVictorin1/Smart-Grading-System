package main

import (
	"net/http"
)

// Routes function to map URL paths to handlers
func (app *application) routes() {
	// Home page route
	http.HandleFunc("/", app.home)

	// Routes for grade management
	http.HandleFunc("/grade", app.viewGrade)
	http.HandleFunc("/create_grade", app.createGrade)
	http.HandleFunc("/grades/edit", app.editGrade)
	http.HandleFunc("/grades/delete", app.deleteGrade)

}
