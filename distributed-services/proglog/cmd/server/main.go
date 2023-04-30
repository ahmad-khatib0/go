package main

import (
	"log"

	"github.com/Ahmadkhatib0/go/distributed-services/proglog/internal/server"
)

func main() {

	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
