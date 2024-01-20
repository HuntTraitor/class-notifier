package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/hunttraitor/class-notifier/internal/models"
	"github.com/hunttraitor/class-notifier/ui"
)

type templateData struct {
	Class           models.Class
	Classes         []models.Class
	Notifications   []models.Notification
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func isFlashError(flash string) bool {
	return flashErrors[flash]
}

var functions = template.FuncMap{
	"isFlashError": isFlashError,
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	//store embed html pages into pages
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		//filepath patterns we want to parse
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		//parse the template files into ts
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		//store into in memory cache
		cache[name] = ts

	}
	return cache, nil
}
