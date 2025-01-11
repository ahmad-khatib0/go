package main

import (
	"fmt"
	"time"
)

func listenToChan(ch chan int) {

	for {
		i := <-ch
		fmt.Println("Got", i, " from channel")

		time.Sleep(time.Second * 1)
	}
}

//  +----------------------------------------------------------------------------------------------------+
//  | what possible value is a buffer channel? And the answer is pretty simple. They're useful when you  |
//  | know how many go routines you've launched. In our case, we've launched one. Or we want to limit    |
//  | the number of go routines we launch, or we want to limit the amount of work that's queued up.      |
//  +----------------------------------------------------------------------------------------------------+

func main() {
	ch := make(chan int, 10)

	go listenToChan(ch)

	for i := 0; i < 100; i++ {
		fmt.Println("sending", i, " to channel...")
		ch <- i
		fmt.Println("sent", i, " to channel...")
	}

	fmt.Println("Done")
	close(ch)
}
