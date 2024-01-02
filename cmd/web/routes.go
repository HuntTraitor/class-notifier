package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/class/add", app.addClass)
	mux.HandleFunc("/class/view", app.viewClass)

	mux.HandleFunc("/notification/add", app.addNotification)
	mux.HandleFunc("/notification/delete", app.deleteNotification)
	mux.HandleFunc("/notification/view", app.viewNotifications)

	return mux
}
