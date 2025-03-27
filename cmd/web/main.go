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

	// read env variables to determine the database connection details
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	unixSocketPath := os.Getenv("INSTANCE_UNIX_SOCKET")

	// this would only be required and provided in Cloud Run deployment
	if unixSocketPath == "" {
		logger.Info(fmt.Sprintf("unixSocketPath: %s", unixSocketPath))
	}

	if dbUser == "" || dbPass == "" || dbName == "" {
		logger.Info(fmt.Sprintf("dbUser: %s", dbUser))
		logger.Info(fmt.Sprintf("dbPass: %s", dbPass))
		logger.Info(fmt.Sprintf("dbName: %s", dbName))
		logger.Error("missing one or more required environment variables for SQL connection")
		os.Exit(1)
	}

	// for local testing - if we did not provide unixSocketPath
	var dsn *string
	if unixSocketPath == "" {
		dsn = flag.String("dsn", dbUser+":"+dbPass+"@tcp(mysql:3306)/"+dbName+"?parseTime=true", "MySQL data source name")
	} else { // for google cloud run - if we provided unixSocketPath
		dsn = flag.String("dsn", dbUser+":"+dbPass+"@unix("+unixSocketPath+")/"+dbName+"?parseTime=true", "MySQL data source name")
	}
	logger.Info(*dsn)

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
