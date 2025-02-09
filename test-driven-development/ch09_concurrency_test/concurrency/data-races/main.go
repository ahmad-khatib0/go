package main

import (
	"fmt"
	"sync"
)

var greetings []string

const workerCount = 3

func greet(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	g := fmt.Sprintf("Hello, friend! I'm Goroutine %d.", id)
	greetings = append(greetings, g)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go greet(i, &wg)
	}

	wg.Wait()
	for _, g := range greetings {
		fmt.Println(g)
	}

	fmt.Println("Goodbye, friend!")
}

// 1- The first data race is detected in the greet function at main.go:15. One goroutine reads
//    a variable, while another goroutine writes to it.
// 2- The second data race is happening as the slice grows during append, which is indicated by the
//    call to runtime.growslice(). This function copies the slice and handles the allocation of a
//    larger backing array, if it is required. The modifications to this slice are also happening in
//    an interleaving manner, with reads and writes happening in different goroutines.
// 3- Finally, the output of the program is printed and the race detector summarizes that two data
//    races have been found.
