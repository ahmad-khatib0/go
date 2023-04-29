package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ahmadkhatib0/go/books-library/internal/data"
	"github.com/Ahmadkhatib0/go/books-library/internal/driver"
)

type config struct {
	port int
}

type application struct {
	config      config
	infoLog     *log.Logger
	errorLog    *log.Logger
	models      data.Models
	environment string
}

func main() {
	var cfg config
	cfg.port = 8082

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN") // this will be sat by: env DSN="the env string" go run ./cmd/api
	// this is the first way, second way is using the makefile (make start)
	environment := os.Getenv("ENV")

	db, err := driver.ConnectPostrges(dsn)
	if err != nil {
		log.Fatal("Can not connect to database")
	}

	defer db.SQL.Close()
	app := &application{
		config:      cfg,
		infoLog:     infoLog,
		errorLog:    errorLog,
		models:      data.New(db.SQL),
		environment: environment,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
