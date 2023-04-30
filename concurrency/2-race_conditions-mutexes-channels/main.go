package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func main() {
	raceCondition()
}

func raceCondition() {
	msg = "hello world"
	wg.Add(2)

	go updateMessage("Hello universe")
	go updateMessage("Hello cosoms")
	wg.Wait()

	fmt.Println(msg)
}
