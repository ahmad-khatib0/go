package main

import "fmt"

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

func main() {
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
