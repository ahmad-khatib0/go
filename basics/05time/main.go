package main

import (
	"fmt"
	"time"
)

func main() {

	presentTime := time.Now()
	fmt.Println(presentTime)

	fmt.Println(presentTime.Format("01-02-2006 15:04:05 Monday")) // NOTE: these numbers are not arbitrary,
	// you need to pass them in this way to get and format the current time

	createdAt := time.Date(2020, time.August, 10, 23, 23, 0, 0, time.UTC)
	fmt.Println(createdAt)
}
