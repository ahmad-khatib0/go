package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
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

	fmt.Println(blankReturns(5, 2)) // 2 1 <nil>

	functionsAreValues()

	anonymousFunctions()

	closures()
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

func blankReturns(numerator, denominator int) (result int, remainder int, err error) {
	// NEVER USE THESE!
	if denominator == 0 {
		err = errors.New("cannot divide by zero")
		return
	}
	result, remainder = numerator/denominator, numerator%denominator
	return // the last line will be returned, but blank return is a bad idea
}

// functions are values
func add(i int, j int) int    { return i + j }
func sub(i int, j int) int    { return i - j }
func mul(i int, j int) int    { return i * j }
func divide(i int, j int) int { return i / j }

type opFuncType func(int, int) int

func functionsAreValues() {
	var opMap = map[string]opFuncType{"+": add, "-": sub, "*": mul, "/": div}
	expressions := [][]string{
		{"2", "+", "3"},
		{"2", "-", "3"},
		{"2", "*", "3"},
		{"2", "/", "3"},
		{"2", "%", "3"},
		{"two", "+", "three"},
		{"5"},
	}
	for _, expression := range expressions {
		if len(expression) != 3 {
			fmt.Println("invalid expression:", expression)
			continue
		}
		p1, err := strconv.Atoi(expression[0])
		if err != nil {
			fmt.Println(err)
			continue
		}
		op := expression[1]
		opFunc, ok := opMap[op]
		if !ok {
			fmt.Println("unsupported operator:", op)
			continue
		}
		p2, err := strconv.Atoi(expression[2])
		if err != nil {
			fmt.Println(err)
			continue
		}
		result := opFunc(p1, p2)
		fmt.Println(result)
	}
}

func anonymousFunctions() {
	for i := 0; i < 5; i++ {
		func(j int) {
			fmt.Println("printing", j, "from inside of an anonymous function")
		}(i)
	}
}

func closures() {
	type Person struct {
		FirstName string
		LastName  string
		Age       int
	}
	people := []Person{
		{"Pat", "Patterson", 37},
		{"Tracy", "Bobbert", 23},
		{"Fred", "Fredson", 18},
	}
	sort.Slice(people, func(i int, j int) bool {
		return people[i].Age < people[j].Age
	})
	fmt.Println(people)

	twoBase := makeMult(2)
	threeBase := makeMult(3)
	for i := 0; i < 3; i++ {
		fmt.Println(twoBase(i), threeBase(i))
	}
}

func makeMult(base int) func(int) int {
	return func(factor int) int {
		return base * factor
	}
}
