package main

import (
	"fmt"
	"time"
)

func main() {

	selectGoRoutine()

	deadLockingGoroutines()

	// infiniteRoutine()

	doneChannel()
}

func process(val int) int {
	return val
}

//  ┏╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍┓
//  ╏ a := <-ch    // reads a value from ch and assigns it to a ╏
//  ╏ ch <- b      // write the value in b to ch                ╏
//  ┗╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍┛

func runThingsConcurrently(in <-chan int, out chan<- int) {
	x := []int{1, 3, 4, 5, 6, 7}
	go func() {
		for val := range x {
			result := process(val)
			out <- result
		}
	}()
}

func selectGoRoutine() {
	myChannel1 := make(chan string)
	myChannel2 := make(chan string)

	go func() {
		myChannel1 <- "dog"
	}()

	go func() {
		myChannel2 <- "cat"
	}()

	// msg := myChannel1
	// fmt.Println(msg)

	select {
	// this main func will be held unloucked until reading from one of either myChannel1 or myChannel2
	case msgFromChannel1 := <-myChannel1:
		fmt.Println(msgFromChannel1)
	case msgFromChannel2 := <-myChannel2:
		fmt.Println(msgFromChannel2)
	}

}

// ********************************* Deadlocking Goroutines *********************************
func deadLockingGoroutines() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		v := 1
		ch1 <- v
		v2 := <-ch2
		fmt.Println(v, v2)
	}()

	v := 2
	var v2 int
	select {
	case ch2 <- v:
	case v2 = <-ch1:
	}
	fmt.Println(v, v2)

}

// ********************************* Buffered VS UnBuffered Channels ************************
//  ┏╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍┓
//  ╏ UnBuffered channels is esentialy used to perform synchourness communication between Goroutines ╏
//  ╏so if the size is zero, or the size is omitted, the channel is unbuffered. (the 2th arg in make)╏
//  ┗╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍╍┛

func bufferedChannels() {
	// a buffered channel has a limited capcity of elements
	myChannel1 := make(chan string, 3)
	chars := []string{"a", "b", "c"}

	for _, s := range chars {
		select {
		case myChannel1 <- s:
		}
	}

	close(myChannel1)

	for result := range myChannel1 { // this indicates that we're able to loop over a closed channel
		fmt.Println(result)
	}

}

func infiniteRoutine() {
	go func() {
		for {
			select {
			default:
				fmt.Println("infinine go routine")
			}
		}
	}()

	time.Sleep(time.Second * 10) // prevent running previous Goroutines for ever
}

// ********************************* The Done Channel *********************************
func doWork(done <-chan bool) {
	// done <- chan   means that this function receive only READ access to the done channel
	for {
		select {
		case <-done:
			// this means that the parent Goroutines (which is go doWork(done) ) the power to
			// cancel this goroutine when it needs to do that, this will prevent infinine unintentionaly goroutine
			return
		default:
			fmt.Println("infinine go routine")
		}
	}
}
func doneChannel() {
	done := make(chan bool)
	go doWork(done)

	time.Sleep(time.Second * 3)
	close(done)
}
