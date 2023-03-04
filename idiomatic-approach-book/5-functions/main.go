package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	result := div(7, 4)
	fmt.Println(result)

	// Simulating Named and Optional Parameters (go doesn't have named or optional params)
	namedAndOptionalParams(MyFuncOpts{LastName: "Noor", Age: 44})

	// variadic params
	fmt.Println(variadicParameters(3, 2, 4, 6, 8))              // [5 7 9 11]
	fmt.Println(variadicParameters(3, []int{1, 2, 3, 4, 5}...)) // [4 5 6 7 8]

	// Multiple Return Values
	res, reminder, err := divAndRemainder(5, 2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(res, reminder) // 2 1

	// Named Return Values
	x, y, z := namedReturn(5, 2)
	fmt.Println(x, y, z) // the name that’s used for a
	// named returned value is local to the function; they don’t enforce any name outside of the function

	fmt.Println(shadowing(5, 2)) // 2 1 nil
	// The values from the return statement were returned even though they
	// were never assigned to the named return parameters
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

func divAndRemainder(numerator int, denominator int) (int, int, error) {
	if denominator == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return numerator / denominator, numerator % denominator, nil
}

func namedReturn(numerator int, denominator int) (result int, remainder int, err error) {
	if denominator == 0 {
		err = errors.New("cannot divide by zero")
		return result, remainder, err
	}
	result, remainder = numerator/denominator, numerator%denominator
	return result, remainder, err
}

func shadowing(numerator, denominator int) (result int, remainder int, err error) {
	// you can shadow a named return value. Be sure that you are assigning to the return value and not to a shadow of it.
	result, remainder = 20, 30 // shadowing by mistake, but didn't affect the accuracy of the return
	if denominator == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return numerator / denominator, numerator % denominator, nil
}
