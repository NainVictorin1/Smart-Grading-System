package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	// Static file server
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Home page
	mux.HandleFunc("GET /{$}", app.home)

	// Grade routes
	mux.HandleFunc("GET /grades", app.viewGrades)          // view all grades
	mux.HandleFunc("POST /grades/submit", app.submitGrade) // submit new grade
	mux.HandleFunc("POST /grades/delete", app.deleteGrade) // delete grade
	mux.HandleFunc("POST /grades/update", app.updateGrade) // update grade

	// Attendance routes
	mux.HandleFunc("GET /attendance", app.viewAttendance)           // view all attendance
	mux.HandleFunc("POST /attendance/mark", app.markAttendance)     // mark present/absent
	mux.HandleFunc("POST /attendance/delete", app.deleteAttendance) // delete attendance
	mux.HandleFunc("POST /attendance/update", app.updateAttendance) // update attendance

	return app.loggingMiddleware(mux)
}
