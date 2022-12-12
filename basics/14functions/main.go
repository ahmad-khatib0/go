package main

import "fmt"

func main() {
	fmt.Println("functions in golang")
	greeter()

	// func ()  { } nested functions are not allowed
	result := add(3, 4)
	fmt.Println("result is ", result)

	firstTypeResult, secondTypeResult := unlimitedAdd(4, 5, 6, 7, 8, 8, 54)
	fmt.Println("unlimitedResult is: ", firstTypeResult, secondTypeResult)
}

func greeter() {
	fmt.Println("hello world")
}

func add(val1 int, val2 int) int {
	return val1 + val2
}

func unlimitedAdd(values ...int) (int, string) {
	total := 0
	for _, val := range values {
		total += val
	}
	return total, " return tow types (int , string) from a function "
}
