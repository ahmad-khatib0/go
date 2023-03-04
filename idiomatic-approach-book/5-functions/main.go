package main

import "fmt"

func main() {
	result := div(7, 4)
	fmt.Println(result)

	// Simulating Named and Optional Parameters (go doesn't have named or optional params)
	namedAndOptionalParams(MyFuncOpts{LastName: "Noor", Age: 44})
}

func div(nominator, denominator int) int {
	// nominator and denominator are of type int
	if denominator == 0 {
		return 0
	}
	return nominator / denominator

}

type MyFuncOpts struct {
	FirstName string
	LastName  string
	Age       int
}

func namedAndOptionalParams(opts MyFuncOpts) error {
	return nil
}
