package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jbh14/nbaoverunders/internal/models"
)

type application struct {
	logger        *slog.Logger
	entries       *models.EntryModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")

	// create a new logger earlier, so that we can log any errors that occur during the startup process (including reading env variables)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// was using hard-coded credentials during testing/dev, but replace with env vars
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbName == "" {
		logger.Error(fmt.Sprintf("dbUser: %s", dbUser))
		logger.Error(fmt.Sprintf("dbPass: %s", dbPass))
		logger.Error(fmt.Sprintf("dbName: %s", dbName))
		logger.Error("missing required environment variables")
		os.Exit(1)
	}

	//sqlWebUserPassword := "R3dmountain"
	//sqlWebUserPassword := "mypassword"
	dsn := flag.String("dsn", dbUser+":"+dbPass+"@tcp(mysql:3306)/"+dbName+"?parseTime=true", "MySQL data source name")

	// building the DSN string directly without command line flag
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)

	// parse the command-line flags
	flag.Parse()

	// create sql.DB connection pool
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close() // defer call so that connection pool is closed before the main() function exits

	// Initialize a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		entries:       &models.EntryModel{DB: db},
		templateCache: templateCache,
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

// go run ./cmd/web
