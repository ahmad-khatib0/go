package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/ahmad-khatib0/go/snippetbox/pkg/models"
)

// The serverError helper writes an error message and stack trace to the errorLog
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.errorLog.Println(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// make a ‘trial’ render by writing the template into a buffer. If this fails, we can respond to the user with an
	// error message. But if it works, we can then write the contents of the buffer to our http.ResponseWriter.
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) authenticatedUser(r *http.Request) *models.User {
	// with this type casting, If there is a *models.User struct in the request context with the key contextKeyUser,
	// then we know that the request is coming from an authenticated-and-valid user,
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}

	return user
}
