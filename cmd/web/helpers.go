package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/NainVictorin1/smart-grade-system/internal/data"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)
	// deal with the error status
	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

const (
	SessionName    = "smart-grade-system-session" // Name of the session cookie
	SessionUserKey = "UserID"                     // Key used to store/retrieve user ID in session
)

// InvalidCSRFHandler responds with 403 Forbidden for invalid or missing CSRF tokens
func (app *application) InvalidCSRFHandler(w http.ResponseWriter, r *http.Request) {
	app.ClientError(w, http.StatusForbidden)
}

// ServerError logs server-side errors and returns a 500 Internal Server Error response
func (app *application) ServerError(w http.ResponseWriter, err error) {
	app.ErrorLog.Output(2, err.Error()) // Log the error with call depth 2 (to report caller)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// ClientError sends a specific status code and its corresponding message to the client
func (app *application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// ContextGetUser extracts the authenticated user from the request context
func (app *application) ContextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value("user").(*data.User)
	if !ok {
		return nil
	}
	return user // Return the authenticated user
}
