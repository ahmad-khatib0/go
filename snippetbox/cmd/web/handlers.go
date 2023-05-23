package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Ahmadkhatib0/go/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't the http.NotFound()
	// function to send a 404 response to the client. becauase you can not changing the catch-all behavior
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

// Add a showSnippet handler function.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "POST" {
	// 	w.Header().Set("Allow", "POST")
	// 	w.Header()["X-XSS-Protection"] = []string{"1; mode=block"} // prevent the name being canonicalized
	// 	w.Header()["Date"] = nil
	// 	// remove header like so.. because Del() method doesn’t remove system-generated headers

	// 	// w.WriteHeader(405)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	// http.Error calls the w.WriteHeader() and .Write() methods behind-the-scenes for us
	// 	return
	// }

	r.Body = http.MaxBytesReader(w, r.Body, 4096) // Limit the request body size to 4096 bytes
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (maximum is 100 characters)"
	}
	// Check that the Content field isn't blank.
	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field cannot be blank"
	}

	if len(errors) > 0 {
		// re-display the create.page.template passing in errors and previously submitted data
		app.render(w, r, "create.page.tmpl", &templateData{
			FormErrors: errors,
			FormData:   r.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
