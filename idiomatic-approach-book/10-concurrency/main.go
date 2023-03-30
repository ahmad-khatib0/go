package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func main() {

	// selectGoRoutine()

	// deadLockingGoroutines()

	// infiniteRoutine()

	// doneChannel()

	// pipelines()

	usingCancelFunctionWithGoroutines()
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

// ********************************      Pipelines     *********************************

func sliceToChannel(nums []int) <-chan int { // read only channel
	out := make(chan int)
	go func() {
		for _, n := range nums {
			fmt.Println("i'm gonna wait 2 sec before writing to the channel")
			time.Sleep(time.Second * 2)
			out <- n
		}
		close(out)
	}()

	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			fmt.Println("i have received the result from the channel")
			out <- n
		}
		close(out)
	}()

	return out // the go routine could be sending data to the out channel even thogh we Returned from the function
}

func pipelines() {
	nums := []int{1, 3, 4, 5, 6} // The Input

	dataChannel := sliceToChannel(nums) // Stage 1

	finalChannel := sq(dataChannel) // Stage 2

	// Stage 3
	for n := range finalChannel {
		fmt.Println(n)
	}
}

// ********************** Goroutines, for Loops, and Varying Variables **********************
func varyingVariables() {
	a := []int{2, 4, 6, 8, 10}
	ch := make(chan int, len(a))
	for _, v := range a {
		go func() {
			ch <- v * 2
		}()
	}
	for i := 0; i < len(a); i++ {
		fmt.Println(<-ch) // 20 20 20 20 20 20
	}

	// The reason why every goroutine wrote 20 to ch is that the closure for every goroutine captured the same
	// variable. The index and value variables in a for loop are reused on each iteration. The last value assigned
	// to v was 10. When the goroutines run, (THAT’S THE VALUE THAT THEY SEE). This problem isn’t unique to for loops;
	// any time a goroutine depends on a variable whose value might change, you must pass the value into the goroutine

	// to solve this problem:  first is to shadow the value within the loop:
	for _, v := range a {
		v := v
		go func() {
			ch <- v * 2
		}()
	}

	// Secondaly: If you want to avoid shadowing and make the data flow more obvious,
	// you can also pass the value as a parameter to the goroutine
	for _, v := range a {
		go func(value int) {
			ch <- value * 2
		}(v)
	}

}

// *********************************  CleanUp Goroutines *********************************
func countTo(max int) <-chan int {
	ch := make(chan int)
	go func() {
		for i := 0; i < max; i++ {
			ch <- i
		}
		close(ch)
	}()
	return ch
}

func cleanupGoroutine() {
	// if we exit the loop early, the goroutine blocks forever, waiting for a value to be read from the channel
	for i := range countTo(10) {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}
}

// Use the done channel to CleanUp the goroutine
// an example, where we pass the same data to multiple functions,
// but only want the result from the fastest function:
func searchData(s string, searchers []func(string) []string) []string {
	// We use an empty struct for the type because the value is unimportant; we never write to this channel, only close it.
	done := make(chan struct{})
	result := make(chan []string)

	for _, searcher := range searchers {
		go func(searcher func(string) []string) {
			select {
			case result <- searcher(s):
			case <-done:
			}
		}(searcher)
	}

	// we read the first value written to result, and then we close done. This signals to the
	// goroutines that they should exit, preventing them from leaking :
	r := <-result
	close(done)
	return r
}

// **********************  Using a Cancel Function To Terminate a Goroutine **********************
func countToWithCancel(max int) (<-chan int, func()) {
	ch := make(chan int)
	done := make(chan struct{})
	cancel := func() {
		close(done)
	}

	go func() {
		for i := 0; i < max; i++ {
			select {
			case <-done:
				return
			default:
				ch <- i
			}
		}
		close(ch)
	}()
	return ch, cancel
}

func usingCancelFunctionWithGoroutines() {
	ch, cancel := countToWithCancel(10)
	for i := range ch {
		if i > 5 {
			break
		}
		fmt.Println(i)
	}
	cancel()
}

// ********************************* Best Way to Use Buffered Channels *******************************
// so we're writing and reading the same amount into and from the results channel, which will ensure
// not having an unintentionaly blocking goroutine that is waiting to be read
// So what we are doing here is reading out the values as they are written. When all of the
// values have been read, we return the results, knowing that we aren’t leaking any goroutines.
func processChannel(ch chan int) []int {
	const conc = 10
	results := make(chan int, conc)
	for i := 0; i < conc; i++ {
		value := i
		go func() {
			results <- process(value)
		}()
	}
	var out []int
	for i := 0; i < conc; i++ {
		out = append(out, <-results)
	}
	return out
}

// ********************************* Backpressure Technique *******************************
// We can use a buffered channel and a select statement to limit the number of simultaneous requests in a system
type PressureGauge struct {
	ch chan struct{}
}

func New(limit int) *PressureGauge {
	ch := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		ch <- struct{}{}
	}

	return &PressureGauge{
		ch: ch,
	}
}

func (pg *PressureGauge) Process(f func()) error {
	select {
	case <-pg.ch:
		f()
		pg.ch <- struct{}{}
		return nil
	default:
		return errors.New("no more capcity")
	}
}

// using it
func doThingThatShouldBeLimited() string {
	time.Sleep(time.Second * 2)
	return "done"
}

func backpressureTechnique() {
	pg := New(10)

	http.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		err := pg.Process(func() {
			w.Write([]byte(doThingThatShouldBeLimited()))
		})
		if err != nil {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests"))
		}
	})
	http.ListenAndServe(":8080", nil)
}

// Turning Off a case In a select
// Assigning the channel to nil, make it unreadable in the next iteration, so the case
// will not execute, so this can prevent consuming precious resources
func avoidJunkCaseInSelect(ch1 <-chan int, ch2 <-chan int, done <-chan bool) {
	for {
		select {
		case _, ok := <-ch1:
			if !ok {
				ch1 = nil // the case will never succeed again!
				continue
			}
		case _, ok := <-ch2:
			if !ok {
				ch2 = nil // the case will never succeed again!
				continue
			}
		case <-done:
			return
		}
	}
}
