package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ahmadkhatib0/go/snippetbox/pkg/models/mysql"
	"github.com/golangcollege/sessions"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	session       *sessions.Session
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	// we use the flag.Parse() function to parse the command-line This reads in the command-line flag
	// value and assigns it to the addr variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000"
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL database connection")

	// Secret for http sessions. Must be 32 bytes long.
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key ")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret)) //  initialize a new session manager,
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		session:       session,
		templateCache: templateCache,
	}

	// if PreferServerCipherSuites true, so the order of the CipherSuites
	// becomes important, (suites which come first in the slice will be preferred.)
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
	}

	// Ports 0-1023 are restricted and (typically) can only be used by services which have root privileges
	infoLog.Printf("Starting server on :4000")
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	// second param needs: Handler , and mux implement a Handler type also, so its staticfied

	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
