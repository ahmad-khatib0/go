package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.AuthTokenMiddleware)
		mux.Post("/users", app.AllUsers)
		mux.Post("/users/save", app.EditUser)
		mux.Post("/users/get/{id}", app.GetUser)
		mux.Post("/users/delete", app.DeleteUser)
		mux.Post("/log-user-out/{id}", app.LogUserOutAndSetInactive)

		mux.Post("/authors/all", app.AuthorsAll)
		mux.Post("/books/save", app.EditBook)
		mux.Post("/books/delete", app.DeleteBook)
		mux.Post("/books/{id}", app.BookById)
	})

	mux.Post("/users/login", app.Login)
	mux.Post("/users/logout", app.Logout)
	mux.Post("/validate-token", app.ValidateToken)

	mux.Post("/books", app.AllBooks)
	mux.Get("/books", app.AllBooks)
	mux.Get("/books/{slug}", app.OneBook)

	// mux.Get("/users/add", func(w http.ResponseWriter, r *http.Request) {
	// 	u := data.User{
	// 		Email:     "you@there.com",
	// 		FirstName: "you",
	// 		LastName:  "there",
	// 		Password:  "password",
	// 	}
	// 	app.infoLog.Println("adding user...")

	// 	id, err := app.models.User.Insert(u)
	// 	if err != nil {
	// 		app.errorLog.Println(err)
	// 		app.errorJSON(w, err, http.StatusForbidden)
	// 		return
	// 	}

	// 	app.infoLog.Println("Got back id of:  ", id)
	// 	newUser, _ := app.models.User.GetOne(id)
	// 	app.writeJSON(w, http.StatusOK, newUser)
	// })

	// mux.Get("/test-generate-token", func(w http.ResponseWriter, r *http.Request) {
	// 	token, err := app.models.Token.GenerateToken(2, 60*time.Minute)
	// 	if err != nil {
	// 		app.errorLog.Println(err)
	// 		return
	// 	}

	// 	token.Email = "you@there.com"
	// 	token.CreatedAt = time.Now()
	// 	token.UpdatedAt = time.Now()

	// 	payload := jsonResponse{
	// 		Error:   false,
	// 		Message: "success",
	// 		Data:    token,
	// 	}

	// 	app.writeJSON(w, http.StatusOK, payload)
	// })

	// mux.Get("/test-save-token", func(w http.ResponseWriter, r *http.Request) {
	// 	token, err := app.models.Token.GenerateToken(7, 60*time.Minute)
	// 	if err != nil {
	// 		app.errorLog.Println(err)
	// 		return
	// 	}

	// 	user, err := app.models.User.GetOne(7)

	// 	token.UserID = user.ID
	// 	token.CreatedAt = time.Now()
	// 	token.UpdatedAt = time.Now()

	// 	err = token.Insert(*token, *user)
	// 	if err != nil {
	// 		app.errorLog.Println(err)
	// 		return
	// 	}

	// 	payload := jsonResponse{
	// 		Error:   false,
	// 		Message: "success",
	// 		Data:    token,
	// 	}

	// 	app.writeJSON(w, http.StatusOK, payload)
	// })

	// mux.Get("/test-validate-token", func(w http.ResponseWriter, r *http.Request) {
	// 	tokenToValidate := r.URL.Query().Get("token")
	// 	valid, err := app.models.Token.ValidToken(tokenToValidate)
	// 	if err != nil {
	// 		app.errorJSON(w, err)
	// 		return
	// 	}

	// 	var payload jsonResponse
	// 	payload.Error = false
	// 	payload.Data = valid

	// 	app.writeJSON(w, http.StatusOK, payload)
	// })

	// static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
