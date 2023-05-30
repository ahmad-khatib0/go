package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	diceNumber := rand.Intn(6) + 1 // between 1 and 6

	switch diceNumber {
	case 1:
		fmt.Println("value is 1")
		fallthrough // matched case 1? will run the case 2 also
	case 2:
		fmt.Println("value is 2")
	case 3:
		fmt.Println("value is 3")
	case 4:
		fmt.Println("value is 4")
	case 5:
		fmt.Println("value is 5")
	case 6:
		fmt.Println("value is 6")
	default:
		fmt.Println("not a number ")
	}
}
