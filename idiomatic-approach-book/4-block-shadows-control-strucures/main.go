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

func loops() {
	// 1- the complete for statement
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	// The Condition-only for Statement
	i := 1
	for i < 100 {
		fmt.Println(i)
		i = i * 2
	}

	// The Infinite for Statement
	// for {
	// 	fmt.Println("Infinite loop")
	// }

	// break and continue
	for i := 1; i <= 100; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("FizzBuzz")
			continue
		}
		if i%3 == 0 {
			fmt.Println("Fizz")
			continue
		}
		if i%5 == 0 {
			fmt.Println("Buzz")
			continue
		}
		fmt.Println(i)
	}

	// The for-range Statement
	evenVals := []int{2, 4, 6, 8, 10, 12}
	for _, v := range evenVals {
		fmt.Println(v)
	}

	uniqueNames := map[string]bool{"Fred": true, "Raul": true, "Wilma": true}
	for k := range uniqueNames {
		fmt.Println(k) // Fred Raul  Wilma
	}

	// Iterating Over Maps
	m := map[string]int{"a": 1, "c": 3, "b": 2}
	for i := 0; i < 3; i++ {
		fmt.Println("Loops", i)
		for k, v := range m {
			fmt.Println(k, v) // The order of the keys and values varies; some runs may be identical.
		}
	}

	// Iterating Over Strings
	samples := []string{"hello", "apple_π!"}
	for _, sample := range samples {
		for i, r := range sample {
			fmt.Println(i, r, string(r)) // convert run to string
		}
		fmt.Println()
	}

	//  Modifying the Value Doesn’t Modify the Source
	for _, v := range evenVals {
		v *= 2
	}
	fmt.Println(evenVals) // [2 4 6 8 10 12] didn't change the value

	// Labeling for nested loops
outer: // note the indentation
	for _, sample := range samples {
		for i, r := range sample {
			fmt.Println(i, r, string(r))
			if r == 'l' {
				continue outer // break the outer loop
			}
		}
	}
}

func main() {
	// shadowingVariables()
	// ifBlock()
	loops()
}
