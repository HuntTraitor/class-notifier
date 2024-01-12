package main

import (
	"net/http"

	"github.com/hunttraitor/class-notifier/ui"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	//serving static files via embed files
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	//adds middleware for dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	//classes handlers with dynamic middleware defined above
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodPost, "/class/add", dynamic.ThenFunc(app.addClass))
	router.Handler(http.MethodGet, "/class/view/:id", dynamic.ThenFunc(app.viewClass))

	//user handlers with dyanimc middleware
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	//appending a new middleware to dynamic for authentication
	protected := dynamic.Append(app.requireAuthentication)

	//using the new protected middleware on top of dynamic for credential specific routes
	router.Handler(http.MethodPost, "/notification/add", protected.ThenFunc(app.addNotification))
	router.Handler(http.MethodPost, "/notification/delete/:id", protected.ThenFunc(app.deleteNotification))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	//Final middleware for all routes
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
