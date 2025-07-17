package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/heisenberg8055/gosts/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

var templateFunctions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2001 at 15:04")
}

func templateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(templateFunctions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
