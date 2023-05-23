package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	// we use the flag.Parse() function to parse the command-line This reads in the command-line flag
	// value and assigns it to the addr variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000"
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Ports 0-1023 are restricted and (typically) can only be used by services which have root privileges
	infoLog.Printf("Starting server on :4000")
	err := srv.ListenAndServe()
	// second param needs: Handler , and mux implement a Handler type also, so its staticfied

	errorLog.Fatal(err)
}
