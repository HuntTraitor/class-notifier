package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hunttraitor/class-notifier/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//Check if the path has any unkown extensions
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	classes, err := app.classes.Classlist()
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.render(w, r, http.StatusOK, "home.html", templateData{
		Classes: classes,
	})
}

// Sends a post request to add the class
func (app *application) addClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
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
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	email := "htratar@ucsc.edu"
	classid := 30262
	expires := 7

	err := app.notifications.Insert(email, classid, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/notification/view", http.StatusSeeOther)

}

func (app *application) deleteNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.Header().Set("Allow", http.MethodDelete)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	email := "htratar@ucsc.edu"
	classid := 31139

	err := app.notifications.Delete(email, classid)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, "/notification/view", http.StatusSeeOther)
}

func (app *application) viewNotifications(w http.ResponseWriter, r *http.Request) {
	classes, err := app.notifications.List("htratar@ucsc.edu")

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, class := range classes {
		fmt.Fprintf(w, "%+v\n", class)
	}
}
