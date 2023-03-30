package main

import (
	"errors"
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
