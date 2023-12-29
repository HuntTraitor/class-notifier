package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/hunttraitor/class-notifier/internal/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// defines application to use for dependency injection for other things
// like middleware and logging
type application struct {
	logger        *slog.Logger
	classes       *models.ClassModel
	notifications *models.NotificationModel
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

	//new instance of an application
	app := &application{
		logger:        logger,
		classes:       &models.ClassModel{DB: db},
		notifications: &models.NotificationModel{DB: db},
	}

	logger.Info("Starting server", "addr", addr)

	err = http.ListenAndServe(addr, app.routes())
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
