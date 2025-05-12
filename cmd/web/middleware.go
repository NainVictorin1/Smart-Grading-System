package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// SecureHeaders sets common HTTP security headers to protect against XSS and clickjacking.
func (app *application) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Protect against cross-site scripting (XSS)
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Prevent content from being embedded in iframes
		w.Header().Set("X-Frame-Options", "deny")
		// Prevent MIME-type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Restrict what is sent in the Referer header
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Log the method, path, and HTTP protocol version
		app.InfoLog.Printf("%s %s %s", r.Method, r.URL.Path, r.Proto)
		// Call the next handler
		next.ServeHTTP(w, r)

		// Log how long the request took
		app.InfoLog.Printf("%s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// EnforceHTTPS redirects HTTP requests to HTTPS and sets HSTS headers.
func (app *application) EnforceHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip for local development
		if r.Host == "localhost:4000" || r.Host == "127.0.0.1:4000" {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("X-Forwarded-Proto") != "https" {
			httpsURL := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, httpsURL, http.StatusPermanentRedirect)
			return
		}

		// Set HTTP Strict Transport Security (HSTS)
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// FlashMessages retrieves flash messages from the session and injects them into the request context.
func (app *application) FlashMessages(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session
		session, err := app.SessionStore.Get(r, SessionName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Check for any flash messages
		if flashes := session.Flashes(); len(flashes) > 0 {
			// Add them to the context
			ctx := context.WithValue(r.Context(), "flashes", flashes)
			r = r.WithContext(ctx)
			// Save the session to remove consumed flashes
			if err := session.Save(r, w); err != nil {
				app.ServerError(w, err)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Authenticate checks for a logged-in user in the session and adds the user to the request context.
func (app *application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve session
		session, err := app.SessionStore.Get(r, SessionName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Extract user ID from session
		userID, ok := session.Values[SessionUserKey].(int)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the user from the database
		user, err := app.users.GetByID(userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Add user to the request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Middleware to check authentication
func (app *application) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a user is in the context (authenticated)
		if app.ContextGetUser(r) == nil {
			session, err := app.SessionStore.Get(r, SessionName)
			if err != nil {
				app.ServerError(w, err)
				return
			}
			// Add a flash message and redirect to login
			session.AddFlash("Please login to access this page")
			if err := session.Save(r, w); err != nil {
				app.ServerError(w, err)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Prevent caching of protected content
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		// Continue to the protected route
		next.ServeHTTP(w, r)
	})
}
