package main

import (
	"net/http"
)

func (app *application) Routes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/home", app.HomeHandler)
	mux.HandleFunc("/login", app.LoginForm)
	mux.HandleFunc("/login/submit", app.LoginHandler)
	mux.HandleFunc("/signup", app.SignupForm)
	mux.HandleFunc("/signup/submit", app.SignupHandler)

	// Authenticated routes (RequireAuthentication middleware applied manually)
	mux.Handle("/user/logout", app.RequireAuthentication(http.HandlerFunc(app.LogoutHandler)))

	// Grade-related routes (require authentication)
	mux.Handle("/grade", app.RequireAuthentication(http.HandlerFunc(app.viewGrade)))
	mux.Handle("/create_grade", app.RequireAuthentication(http.HandlerFunc(app.createGrade)))
	mux.Handle("/grades/edit", app.RequireAuthentication(http.HandlerFunc(app.editGrade)))
	mux.Handle("/grades/delete", app.RequireAuthentication(http.HandlerFunc(app.deleteGrade)))
	// Static file server for serving assets (e.g., CSS, JS, images)
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Return the router wrapped with standard middleware
	return app.recoverPanicMiddleware(
		app.SecureHeaders(
			app.logRequestMiddleware(
				app.EnforceHTTPS(
					app.FlashMessages(
						app.Authenticate(mux),
					),
				),
			),
		),
	)
}
