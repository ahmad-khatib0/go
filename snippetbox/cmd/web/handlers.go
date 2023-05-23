package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't the http.NotFound()
	// function to send a 404 response to the client. becauase you can not changing the catch-all behavior
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// NOTE: that the home.page.tmpl file must be the *first* file in the slice.
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	//  store the templates in a template set.

	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}

}

// Add a showSnippet handler function.
func (app *application) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
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

	w.Write([]byte("Create a new snippet..."))
}
