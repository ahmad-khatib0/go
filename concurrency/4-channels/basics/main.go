package main

import (
	"fmt"
	"strings"
)

// 1th  receive only ,  2th send only channels
func shout(ping <-chan string, pong chan<- string) {

	for {
		s := <-ping

		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
	}

}

func main() {
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong) // will run forever in the background

	fmt.Println("type something and hit Enter (enter Q to quit)")

	for {
		fmt.Println("-> ")

		var userInput string
		_, _ = fmt.Scanln(&userInput) // read what user typed

		if userInput == strings.ToLower("q") {
			break
		}

		ping <- userInput
		response := <-pong // wait for a response

		fmt.Println("response: ", response)
	}

	fmt.Println("All done , closing the channels")
	close(ping)
	close(pong)
}
