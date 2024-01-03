package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/class/add", app.addClass)
	mux.HandleFunc("/class/view", app.viewClass)

	mux.HandleFunc("/notification/add", app.addNotification)
	mux.HandleFunc("/notification/delete", app.deleteNotification)
	mux.HandleFunc("/notification/view", app.viewNotifications)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
