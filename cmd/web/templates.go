package main

import (
	"asniki/snippetbox/internal/models"
	"asniki/snippetbox/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

// templateData holds dynamic data to pass to the HTML templates
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// humanDate returns a nicely formatted string representation of a time.Time object
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04:05")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache initializes a new template cache
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.
			New(name).
			Funcs(functions).
			ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
