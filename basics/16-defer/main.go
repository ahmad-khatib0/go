package main

import "fmt"

func main() {

	defer fmt.Println("world")
	defer fmt.Println("one")
	defer fmt.Println("tow")

	fmt.Println("hello")
	loopDefers() // will run or will be logged out before the previous defers ,

	// defers will be placed just at the end of the function,, multiple defers will execute in the reverse order,
	// or last in first out
}

func loopDefers() {
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
}
