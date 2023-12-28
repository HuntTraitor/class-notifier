package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// defines application to use for dependency injection for other things
// like middleware and logging
type application struct {
	logger *slog.Logger
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

	//new instance of an application
	app := &application{
		logger: logger,
	}

	logger.Info("Starting server", "addr", addr)

	err = http.ListenAndServe(addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
