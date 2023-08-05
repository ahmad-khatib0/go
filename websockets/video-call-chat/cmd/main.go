package main

import (
	"log"

	"github.com/ahmad-khatib0/go/websockets/video-call-chat/internal/server"
)

func main() {

	if err := server.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
