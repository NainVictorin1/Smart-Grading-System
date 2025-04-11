package main

import (
	"fmt"
	"net/http"
	"time"
)

// ServeHTTP initializes and starts the HTTP server.
func (app *application) ServeHTTP() error {
	// Register the routes
	app.routes()

	// Create a custom HTTP server with timeouts
	server := &http.Server{
		Addr:           *app.addr,            // Address for the server to listen on
		Handler:        http.DefaultServeMux, // Default mux to handle registered routes
		ReadTimeout:    10 * time.Second,     // Max duration to read the request
		WriteTimeout:   10 * time.Second,     // Max duration for writing a response
		MaxHeaderBytes: 1 << 20,              // Max header size (1 MB)
	}

	// Log and start the server
	fmt.Printf("Server starting on %s\n", *app.addr)
	return server.ListenAndServe()
}
