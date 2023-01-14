package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

func Test_Routes_Exist(t *testing.T) {
	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	routeExists(t, chiRoutes, "/users/login")
	routeExists(t, chiRoutes, "/users/logout")
	routeExists(t, chiRoutes, "/admin/users/get/{id}")
	routeExists(t, chiRoutes, "/admin/users/save")
	routeExists(t, chiRoutes, "/admin/users")
	routeExists(t, chiRoutes, "/admin/users/delete")
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	// assume that the route does not exist
	found := false
	_ = chi.Walk(
		routes,
		func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			if route == foundRoute {
				found = true
			}
			return nil
		},
	)

	if !found {
		t.Errorf("did not found %s in registered routes", route)
	}
}
