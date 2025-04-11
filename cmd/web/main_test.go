package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.home)
	handler.ServeHTTP(rr, req)

	// Check the status code
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}

	// Check the response body
	expected := "Welcome\n" // Adjust to the actual output in your handler
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
