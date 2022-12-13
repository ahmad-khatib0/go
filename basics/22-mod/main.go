package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("mod in golang")
	greeter()
	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")

	log.Fatal(http.ListenAndServe(":4000", r))
}

func greeter() {
	fmt.Println("hello world")
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<ht> the first use of the routing logic in golang </h1>"))
}
