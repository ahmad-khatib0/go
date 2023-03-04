package main

import (
	"fmt"
)

func shadowingVariables() {
	x := 10
	if x > 5 {
		fmt.Println(x) // 10
		x := 5         //  the shadowing variable
		fmt.Println(x) // 5
	}
	fmt.Println(x) // 10

	// shadowing with multiple assignment
	a := 10
	if a > 5 {
		a, y := 5, 20     // using :=  makes it unclear what variables are being used, and accidentally shadowing variables
		fmt.Println(a, y) // 5 20
	}
	fmt.Println(x) // 10

	// shadowing package names
	b := 10
	fmt.Println(b)
	// fmt := "oops" // accidentally shadowing a package
	// fmt.Println(fmt) // an ERROR

}

func main() {
	shadowingVariables()
}
