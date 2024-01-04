package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hunttraitor/class-notifier/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	classes, err := app.classes.Classlist()
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	notifications, err := app.notifications.NotificationList("htratar@ucsc.edu")
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.render(w, r, http.StatusOK, "home.html", templateData{
		Classes:       classes,
		Notifications: notifications,
	})
}

// Sends a post request to add the class
func (app *application) addClass(w http.ResponseWriter, r *http.Request) {

	classid := 31139
	name := "CSE 102 - 01   Introduction to Analysis of Algorithms"
	link := "https://pisa.ucsc.edu/class_search/index.php?action=detail&class_data=YToyOntzOjU6IjpTVFJNIjtzOjQ6IjIyNDAiO3M6MTA6IjpDTEFTU19OQlIiO3M6NToiMzExMzkiO30%253D"
	professor := "Chatziafratis,E."

	id, err := app.classes.Insert(classid, name, link, professor)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/class/view?id=%d", id), http.StatusSeeOther)
}

// redircts to the link for the class
func (app *application) viewClass(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	class, err := app.classes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", class)
}

func (app *application) addNotification(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	classid, err := strconv.Atoi(r.PostFormValue("classId"))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	email := "htratar@ucsc.edu"
	expires := 7

	err = app.notifications.Insert(email, classid, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *application) deleteNotification(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	notificationid, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.notifications.Delete(notificationid)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) viewNotifications(w http.ResponseWriter, r *http.Request) {

	classes, err := app.notifications.NotificationList("htratar@ucsc.edu")

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, class := range classes {
		fmt.Fprintf(w, "%+v\n", class)
	}
}
