package main

import (
	"crypto/tls"
	"database/sql"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hunttraitor/class-notifier/internal/models"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// defines application to use for dependency injection for other things
// like middleware and logging
type application struct {
	logger         *slog.Logger
	classes        models.ClassModelInterface
	notifications  models.NotificationModelInterface
	users          models.UserModelInterface
	sessionManager *scs.SessionManager
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env file could not be loaded")
	}

	//loads addr from PORT env variable
	addr := ":" + os.Getenv("PORT")
	//new logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(os.Getenv("POSTGRES_DSN"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	formDecoder := form.NewDecoder()

	//new instance of an application and server
	app := &application{
		logger:         logger,
		classes:        &models.ClassModel{DB: db},
		notifications:  &models.NotificationModel{DB: db},
		users:          &models.UserModel{DB: db},
		sessionManager: sessionManager,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
	}

	logger.Info("Starting server", "addr", srv.Addr)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
