package main

import (
	"fmt"
	"math/cmplx"
)

func main() {
	x := complex(2.5, 3.1)
	y := complex(10.2, 2)
	fmt.Println(x + y)        // (12.7+5.1i)
	fmt.Println(x - y)        // (-7.699999999999999+1.1i)
	fmt.Println(x * y)        // (19.3+36.62i)
	fmt.Println(x / y)        // (0.2934098482043688+0.24639022584228065i)
	fmt.Println(real(x))      // (2.5)
	fmt.Println(imag(x))      // (3.1)
	fmt.Println(cmplx.Abs(x)) // (3.982461550347975)

	// Type Conversions
	var a int = 10
	var b float64 = 32.5
	var c float64 = float64(a) + b // not identical types, we need to convert them to add them together.
	var d int = int(b) + a
	fmt.Println(c, d)

	// variables declarations ways
	var e int             // if you want to declare a variable and assign it the zero value
	var f, g int = 10, 20 //  declare multiple variables at once with var
	h, i := 10, "hello"
	fmt.Println(e, f, g, h, i)

	var (
		j    int
		k        = 20
		l    int = 30
		m, n     = 40, "hello"
		o, p string
	)
	fmt.Println(j, k, l, m, n, o, p)
	// short declaration variables The := operator can do one trick that you cannot do with var;
	// it allows you to assign values to existing variables, too. As long as there
	// is one new variable on the left hand side of the :=
	q := 10
	q, s := 30, "hello"
	fmt.Println(q, s)
}

func constants() {
	// untyped constant declaration
	const x = 10 // All of the following assignments are legal:
	var y int = x
	var z float64 = x
	var d byte = x

	const typedX int = 10 // typed constant declaration
	fmt.Println(y, z, d)
}

func unusedVariables() {
	x := 10
	x = 20
	fmt.Println(x)
	x = 30
	// While the compiler and go vet do not catch the unused assignments of 10 and 30 to x, golangci-lint, detects them
}
