package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type application struct {
	Logger *logger.Logger
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Sync()

	a := application{
		Logger: l,
	}

	// cdir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal("failed to get current dir %w", err)
	// }

	executable, _ := os.Executable()
	a.Logger.Info(fmt.Sprintf("server %s", executable))

}
