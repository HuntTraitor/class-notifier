package main

import (
	"html/template"
	"path/filepath"

	"github.com/hunttraitor/class-notifier/internal/models"
)

type templateData struct {
	Class           models.Class
	Classes         []models.Class
	Notifications   []models.Notification
	Form            any
	Flash           string
	IsAuthenticated bool
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		//parse base template files into ts
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		//parse all partials into ts
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		//call parsefiles on ts
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		//add ts to map for caching
		cache[name] = ts
	}
	return cache, nil
}
