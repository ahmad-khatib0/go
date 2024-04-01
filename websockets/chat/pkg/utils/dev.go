package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
)

func (u *Utils) LogConfigToFile() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("json.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(enc)
	if err != nil {
		log.Fatal(err)
	}
}
