package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"os"
)

func defers() {
	defer func() int {
		return 2 // there's no way to read this value
	}()

	if len(os.Args) < 2 {
		log.Fatal("no file specified")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	data := make([]byte, 2048)
	for {
		count, err := f.Read(data)
		os.Stdout.Write(data[:count])
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
	}

	// InsertIntoDb(context.Background() , sql.Open() , 1, 1)

	_, closer, err := getFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer closer()
}

func InsertIntoDb(ctx context.Context, db *sql.DB, value1, value2 int) (err error) {
	// an example to handle database transaction cleanup using named return values and defer
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	if err != nil {
		return err
	}
	// use tx to do more database inserts here
	return nil
}

func getFile(name string) (*os.File, func(), error) {
	// A common pattern in Go is for a function that allocates a resource to also return
	// a closure that cleans up the resource
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	return file, func() {
		file.Close()
	}, err
}
