package main

import (

	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"github.com/jbh14/nbaoverunders/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
	entries *models.EntryModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	sqlWebUserPassword := "R3dmountain"
	dsn := flag.String("dsn", "web:" + sqlWebUserPassword + "@/nbaoverunders?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// create sql.DB connection pool
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()  // defer call so that connection pool is closed before the main() function exits

	app := &application{
		logger: 	logger,
		entries: 	&models.EntryModel{DB: db},
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}