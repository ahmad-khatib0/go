package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ahmad-khatib0/go/test-driven-development/ch08_microservices_test/db"
	"github.com/ahmad-khatib0/go/test-driven-development/ch08_microservices_test/handlers"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	port, ok := os.LookupEnv("BOOKSWAP_PORT")
	if !ok {
		log.Fatal("$BOOKSWAP_PORT not found")
	}

	postgresURL, ok := os.LookupEnv("BOOKSWAP_DB_URL")
	if !ok {
		log.Fatal("env variable BOOKSWAP_DB_URL not found")
	}

	m, err := migrate.New("file://chapter08/db/migrations", postgresURL)
	if err != nil {
		log.Fatalf("migrate:%v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration up:%v", err)
	}
	// defer func() {
	// 	m.Down()
	// }()

	dbConn, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("db open:%v", err)
	}

	ps := db.NewPostingService()
	b := db.NewBookService(dbConn, ps)
	u := db.NewUserService(dbConn, b)
	h := handlers.NewHandler(b, u)

	router := handlers.ConfigureServer(h)
	log.Printf("Listening on :%s...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", port), router))
}
