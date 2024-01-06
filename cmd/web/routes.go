package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodPost, "/class/add", dynamic.ThenFunc(app.addClass))
	router.Handler(http.MethodGet, "/class/view/:id", dynamic.ThenFunc(app.viewClass))

	router.Handler(http.MethodPost, "/notification/add", dynamic.ThenFunc(app.addNotification))
	router.Handler(http.MethodPost, "/notification/delete/:id", dynamic.ThenFunc(app.deleteNotification))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
