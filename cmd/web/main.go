package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/NainVictorin1/smart-grade-system/internal/data"

	_ "github.com/lib/pq"
)

type application struct {
	addr          *string
	feedback      *data.FeedbackModel
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	dsn := flag.String("dsn", "", "PostgreSQL DSN")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database connection pool established")
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		addr:          addr,
		feedback:      &data.FeedbackModel{DB: db},
		logger:        logger,
		templateCache: templateCache,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
