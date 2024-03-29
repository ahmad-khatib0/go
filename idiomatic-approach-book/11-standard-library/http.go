package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// httpRequest()

	middleware()
}

// ***************** HTTP Client
func httpRequest() {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req.Header.Add("X-My-Client", "Learning Go")

	res, err := client.Do(req) // send it
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status: got %v", res.Status))
	}

	fmt.Println(res.Header.Get("Content-Type")) // application/json; charset=utf-8
	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", data) // { UserID:1 ID:1 Title:delectus aut autem Completed:false }
}

// ***************** HTTP Server
type HelloHandler struct{}

func (hh HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!\n"))

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      HelloHandler{},
	}
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func muxServe() {

	person := http.NewServeMux()
	person.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("greetings!\n"))
	})

	dog := http.NewServeMux()
	dog.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("good puppy!\n"))
	})

	mux := http.NewServeMux()
	mux.Handle("/person/", http.StripPrefix("/person", person))
	mux.Handle("/dog/", http.StripPrefix("/dog", dog))
}

// *****************  Middleware

func RequestTimer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		end := time.Now()
		log.Printf("request time for %s: %v", r.URL.Path,
			end.Sub(start))
	})
}

var securityMsg = []byte("You didn't give the secret password\n")

func TerribleSecurityProvider(password string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if r.Header.Get("X-Secret-Password") != password {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(securityMsg)
				return
			}

			h.ServeHTTP(w, r)

		})
	}
}

func middleware() {
	mux := http.NewServeMux()

	terribleSecurity := TerribleSecurityProvider("GOPHER")

	mux.Handle("/hello", terribleSecurity(RequestTimer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello!\n"))
		}),
	)),
	)

	// you can apply a set of middleware to all of the handlers registered with a single request router:
	terribleSecurity = TerribleSecurityProvider("GOPHER")
	wrappedMux := terribleSecurity(RequestTimer(mux))
	s := http.Server{
		Addr: ":8080",

		Handler: wrappedMux,
	}

	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}

}
