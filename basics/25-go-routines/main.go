package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	wg  sync.WaitGroup // pointer
	mut sync.Mutex     // pointer
)

var signals = []string{"test"} // pointer

func main() {
	// go greeter("hello")
	// greeter("world")
	websiteList := []string{
		"https://lco.dev",
		"https://go.dev",
		"https://google.com",
		"https://fb.com",
		"https://github.com",
	}

	for _, web := range websiteList {
		go getStatusCode(web)
		wg.Add(1) // 1 because its one job
	}

	wg.Wait() // will make sure that main func will not be ended before finishing the goroutines
	fmt.Println(signals)
}

// func greeter(s string) {
// 	for i := 0; i < 6; i++ {
// 		time.Sleep(3 * time.Millisecond)
// 		fmt.Println(s)
// 	}
// }

func getStatusCode(endpoint string) {
	defer wg.Done() // your responsibility to call this done method

	res, err := http.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	} else {
		mut.Lock()
		signals = append(signals, endpoint)
		mut.Unlock()

		fmt.Printf("%d status code for the website %s\n", res.StatusCode, endpoint)
	}
}
