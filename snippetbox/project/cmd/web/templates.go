package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ahmad-khatib0/go/snippetbox/project/pkg/forms"
	"github.com/ahmad-khatib0/go/snippetbox/project/pkg/models"
	"github.com/justinas/nosurf"
)

type templateData struct {
	AuthenticatedUser *models.User
	CSRFToken         string
	CurrentYear       int
	Flash             string
	Form              *forms.Form
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl')
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before call the ParseFiles() method. This
		// means we have to use template.New create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}
	return cache, nil

}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash") // Add the flash message to the template data, if one exists.

	return td
}

// Create a humanDate function which returns a nicely formatted string representation of a time.Time object.
func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
