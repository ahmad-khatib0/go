package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ahmad-khatib0/go/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		// they essentially instruct the user’s web browser to implement some additional security measures to help prevent
		// XSS and Clickjacking attacks. It’s good practice to include them unless you have a specific reason for not doing so.

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or not. If there has...
			if err := recover(); err != nil {
				// The value returned by the builtin recover() function is an interface{}
				// and its underlying type could be string, error, or something else

				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")

				// Call the app.serverError helper method to return a 500 Internal Server response.
				app.serverError(w, fmt.Errorf("%s", err))

			}
		}()

		next.ServeHTTP(w, r)

	})

}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.authenticatedUser(r) == nil {
			http.Redirect(w, r, "/user/login", 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("called authenticate")
		// When we don’t have an authenticated-and-valid user, we pass the original and unchanged
		// *http.Request to  the next handler in the chain. (notice also the next check)
		exists := app.session.Exists(r, "userID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.Get(app.session.GetInt(r, "userID"))
		if err == models.ErrNoRecord {
			app.session.Remove(r, "userID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		// We create a new copy of the request with the user information added to the request context, and
		// call the next handler in the chain *using this new copy of the request*.
		ctx := context.WithValue(r.Context(), contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
