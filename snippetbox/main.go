package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't the http.NotFound()
	// function to send a 404 response to the client. becauase you can not changing the catch-all behavior
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox"))
}

// Add a showSnippet handler function.
func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a createSnippet handler function.
func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		w.Header().Set("Allow", "POST")

		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"} // prevent the name being canonicalized
		w.Header()["Date"] = nil
		// remove header like so.. because Del() method doesn’t remove system-generated headers

		// w.WriteHeader(405)
		http.Error(w, "Method Not Allowed", 405)
		// http.Error calls the w.WriteHeader() and .Write() methods behind-the-scenes for us
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
