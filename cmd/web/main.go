package main

import (
	"asniki/snippetbox/internal/models"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	// import is needed for the driver’s init() function to run so that it can register itself with the database/sql package
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	addr      string
	staticDir string
	dsn       string
	app       *application
	tlsConfig *tls.Config
	logLogger *log.Logger
)

// application holds the application-wide dependencies for the web application
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

// init parses cmd flags initializes application struct and other variables for http.Server
func init() {
	logger := initLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Error(fmt.Sprintf("Error loading .env file: %v", err.Error()))
		os.Exit(1)
	}
	defaultDsn := os.Getenv("DSN")

	flag.StringVar(&addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&dsn, "dsn", defaultDsn, "MySQL data source name")
	flag.Parse()

	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app = &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig = &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
}

// initLogger initializes a new structured logger
func initLogger() *slog.Logger {
	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					source.File = filepath.Base(source.File)
				}
			}
			return a
		},
	})
	logLogger = slog.NewLogLogger(loggerHandler, slog.LevelError)
	return slog.New(loggerHandler)
}

// openDB initializes DB connection pool and check connection for errors
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

func main() {
	app.logger.Info("Starting server on",
		"addr", addr)

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     logLogger,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.logger.Error(err.Error())
	os.Exit(1)
}
