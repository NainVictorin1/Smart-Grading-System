package main

import (
	"html/template"
	"path/filepath"
)

// newTemplateCache loads all HTML templates into a map for fast access.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Update the pattern to search for .tmpl files in the directory
	// Make sure we match the templates
	pages, err := filepath.Glob("./ui/html/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		fileName := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[fileName] = ts
	}

	return cache, nil
}
