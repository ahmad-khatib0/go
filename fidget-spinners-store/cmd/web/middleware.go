package main

import "net/http"

// we have access to session because that's a package level variable

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
