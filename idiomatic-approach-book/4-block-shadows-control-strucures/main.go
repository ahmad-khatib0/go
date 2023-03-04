package main

import (
	"fmt"
	"math/rand"
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

func ifBlock() {
	// Go adds is the ability to declare variables that are scoped to the condition and to
	// both the if and else blocks. here Once the series of if/else statements ends, n is undefined.
	if n := rand.Intn(10); n == 0 {
		fmt.Println("That's too low")
	} else if n > 5 {
		fmt.Println("That's too big:", n)
	} else {
		fmt.Println("That's a good number:", n)
	}
}

func main() {
	// shadowingVariables()
	ifBlock()
}
