package main

import (
	"fmt"
	"sync"
)

// channels are a way in which multiple go routines can talk to each other,

func main() {
	fmt.Println("channels in golang")

	wg := &sync.WaitGroup{}
	channel := make(chan int, 2)

	// channel <- 4 // push
	// fmt.Println(<-channel) // throws an error, the summary of the problem with channels
	//  is like so : i don't allow you to push a value if no one is listening to me

	wg.Add(2)
	go func(ch <-chan int, wg *sync.WaitGroup) {
		// this <-chan to prevent putting e.g close(channel) unintentionally before reading from channel
		// this this <-chan indicate that this function here will only listen to channel, AND WILL NOT MODIFY IT

		// this if block is to handle the NOTE issue
		val, isChannelOpen := <-channel
		if isChannelOpen {
			fmt.Println(val)
		}

		// fmt.Println(<-channel)
		// fmt.Println(<-channel)
		wg.Done()
	}(channel, wg)
	go func(ch chan<- int, wg *sync.WaitGroup) {
		channel <- 4
		channel <- 8 // if we made the channel chan as un-buffered channel, so listening on tow cases
		// will throws an error, because in the first func where we listen, we listen only for the first push

		// NOTE: jf you close the channel without pushing anything, it will return zero to you,
		// and if you pushed zero to the channel and than closed the channel, it also return zero,
		// so how to differentiate if the zero is pushed by us, or its the one that is returned be default ?

		close(channel)
		wg.Done()
	}(channel, wg)

	wg.Wait()
}
