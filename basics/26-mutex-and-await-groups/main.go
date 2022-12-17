package main

import (
	"fmt"
	"sync"
)

// race condetion is happening when you are writting to the same memory location by tow or more
// goroutines at the same time, which will cause many issues, or which is shouldn't happening
func main() {
	fmt.Println("Race Condetion in golang")

	wg := &sync.WaitGroup{}
	mut := &sync.RWMutex{}

	// mut.RLock() // don't lock this resource here, instead lock it WHEN you are READING from it (4th func)
	score := []int{0}
	// mut.RUnlock()

	wg.Add(4)
	// wg.Add(1) // intstead of adding one each time for each function
	go func(wg *sync.WaitGroup, m *sync.RWMutex) {
		fmt.Println("routine number 1")
		// so this is saying whenever we work on memory , just lock it and don't allow other job to use it
		mut.Lock()
		score = append(score, 1)
		mut.Unlock()
		wg.Done()
	}(wg, mut)

	// wg.Add(1)
	go func(wg *sync.WaitGroup, m *sync.RWMutex) {
		fmt.Println("routine number 2")
		mut.Lock()
		score = append(score, 2)
		mut.Unlock()
		wg.Done()
	}(wg, mut)

	// wg.Add(1)
	go func(wg *sync.WaitGroup, m *sync.RWMutex) {
		fmt.Println("routine number 3")
		mut.Lock()
		score = append(score, 3)
		mut.Unlock()
		wg.Done()
	}(wg, mut)

	go func(wg *sync.WaitGroup, m *sync.RWMutex) {
		fmt.Println("routine number 4")
		mut.RLock()
		fmt.Println(score)
		mut.RUnlock()
		wg.Done()
	}(wg, mut)

	wg.Wait()
	fmt.Println(score)
}
