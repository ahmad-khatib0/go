package main

import (
	"fmt"
	"sync"
)

func printSomething(s string, wg *sync.WaitGroup) {
	defer wg.Done() // Done to decrease the number of 9 on each item by one
	fmt.Println(s)
}

func waitGroups() {
	fmt.Println("hello concurrency")
	var wg sync.WaitGroup

	// this slice consists of the words we want to print using a goroutine
	words := []string{
		"alpha",
		"beta",
		"delta",
		"gamma",
		"pi",
		"zeta",
		"eta",
		"theta",
		"epsilon",
	}

	wg.Add(len(words))

	for i, x := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, x), &wg)
	}

	wg.Wait()

	wg.Add(1)
	printSomething("this is the second thing to be printed!", &wg)

}

func main() {
	// waitGroups()
	Challenge()
}
