package main

import (
	"errors"
	"sync"
	"time"
)

func main() {

}

// Time Out Code

func timeLimit() (int, error) {

	var result int
	var err error
	done := make(chan struct{})
	go func() {
		result, err = doSomeWork()
		close(done)
	}()

	select {
	case <-done:
		return result, err
	case <-time.After(2 * time.Second):
		return 0, errors.New("Sorry, your request timed out")
	}
}

func doSomeWork() (int, error) {
	time.Sleep(time.Second * 1)
	return 2, nil
}

// *********************************   Using WaitGroups  *********************************
func waitGroups() {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		doSomeWork1()
	}()
	go func() {
		defer wg.Done()
		doSomeWork2()
	}()
	go func() {
		defer wg.Done()
		doSomeWork3()
	}()

	wg.Wait()
}

func doSomeWork1() {}
func doSomeWork2() {}
func doSomeWork3() {}

// *********************************   Using WaitGroups  *********************************
func processAndGather(in <-chan int, processor func(int) int, num int) []int {
	out := make(chan int, num)
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for v := range in {
				//  The for-range channel loop exits when out is closed and the buffer is empty
				out <- processor(v)
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	var result []int
	for v := range out {
		result = append(result, v)
	}
	return result
}
