package main

import "fmt"

func main() {
	result := div(7, 4)
	fmt.Println(result)

	// Simulating Named and Optional Parameters (go doesn't have named or optional params)
	namedAndOptionalParams(MyFuncOpts{LastName: "Noor", Age: 44})

	fmt.Println(variadicParameters(3, 2, 4, 6, 8))              // [5 7 9 11]
	fmt.Println(variadicParameters(3, []int{1, 2, 3, 4, 5}...)) // [4 5 6 7 8]
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

func variadicParameters(base int, vals ...int) []int {
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		out = append(out, base+v)
	}
	return out
}
