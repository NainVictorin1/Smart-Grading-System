package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/NainVictorin1/smart-grade-system/internal/data"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

type application struct {
	addr          *string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	logger        *slog.Logger
	DB            *sql.DB
	grades        *data.GradeModel
	users         *data.UserModel
	templateCache map[string]*template.Template
	CSRFKey       []byte
	SessionStore  *sessions.CookieStore
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address") // Set a default address
	dsn := flag.String("dsn", "postgres://teacher:nain@localhost/grade?sslmode=disable", "PostgreSQL DSN")
	sessionKey := flag.String("session-key", "Zs6yBsEyTRu/Hw5x/tw2tSmR1VJEeCPKCdV88WU0gR8=", "Session encryption key")
	csrfKey := flag.String("csrf-key", "hD6VrOk/pCu8F7DWGNBHvbShSXZDC8W+jc4z/XBuwIY=", "CSRF encryption key")
	flag.Parse()

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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
	sessionStore := sessions.NewCookieStore([]byte(*sessionKey))
	sessionStore.Options = &sessions.Options{
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400 * 7,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Set to true in production with HTTPS
	}
	// Initialize application with the logger, grade model, and template cache
	app := &application{
		addr:          addr,
		SessionStore:  sessionStore,
		users:         &data.UserModel{DB: Database},
		grades:        &data.GradeModel{DB: Database},
		logger:        logger,
		templateCache: templateCache,
		CSRFKey:       []byte(*csrfKey),
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
	}
	// Set up CSRF protection middleware
	csrfMiddleware := csrf.Protect(
		app.CSRFKey,
		csrf.Secure(false),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteLaxMode), // Protects against some types of CSRF attacks
		csrf.HttpOnly(true),
		csrf.FieldName("csrf_token"),
		csrf.ErrorHandler(http.HandlerFunc(app.InvalidCSRFHandler)),
	)

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      csrfMiddleware(app.Routes()), // Use the implemented routes() function
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlsConfig,
	}
	app.logger.Info("Starting server on port", "port", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		app.logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
