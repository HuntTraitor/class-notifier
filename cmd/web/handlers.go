package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//Check if the path has any unkown extensions
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	//add all relevant files to slice
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	//Read the template file into ts
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Execute the template to write in responsewriter
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, r, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Sends a post request to add the class
func (app *application) addClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Add a class for a user..."))
}

// redircts to the link for the class
func (app *application) viewClass(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "Display a class with name %s...", name)
}
