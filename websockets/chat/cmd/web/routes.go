package main

import (
	"net/http"
	"ws/internal/handlers"

	"github.com/bmizerany/pat"
)

// routes defines the application routes
func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))

	return mux
}
