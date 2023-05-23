package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Ahmadkhatib0/go/snippetbox/pkg/models"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't the http.NotFound()
	// function to send a 404 response to the client. becauase you can not changing the catch-all behavior
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// NOTE: that the home.page.tmpl file must be the *first* file in the slice.
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)//  store the templates in a template set.
	// err = ts.Execute(w, nil)

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}
}

// Add a showSnippet handler function.
func (app *application) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%v", s)
}

// Add a createSnippet handler function.
func (app *application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		w.Header().Set("Allow", "POST")

		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"} // prevent the name being canonicalized
		w.Header()["Date"] = nil
		// remove header like so.. because Del() method doesn’t remove system-generated headers

		// w.WriteHeader(405)
		app.clientError(w, http.StatusMethodNotAllowed)
		// http.Error calls the w.WriteHeader() and .Write() methods behind-the-scenes for us
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
