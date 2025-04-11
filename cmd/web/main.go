package main

import (
	"flag"
	"html/template"
	"log/slog"
	"os"

	"github.com/NainVictorin1/smart-grade-system/internal/data"
)

type application struct {
	addr          *string
	logger        *slog.Logger
	grades        *data.GradeModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address") // Set a default address
	dsn := flag.String("dsn", "", "PostgreSQL DSN")

	flag.Parse()

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize the database connection
	err := initDatabase(*dsn)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}

	logger.Info("Database connection pool established")

	// Load the template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error("Failed to load templates", "error", err)
		os.Exit(1)
	}

	defer Database.Close()

	// Initialize application with the logger, grade model, and template cache
	app := &application{
		addr:          addr,
		grades:        &data.GradeModel{DB: Database},
		logger:        logger,
		templateCache: templateCache,
	}

	// Start the HTTP server using the ServeHTTP method from server.go
	err = app.ServeHTTP()
	if err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
