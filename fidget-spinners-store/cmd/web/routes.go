package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// mux.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"https://*", "http://*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300,
	// }))

	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)

	return mux
}
