package main

import "fmt"

func main() {
	var x int32 = 10
	var y bool = true
	pointerX := &x
	pointerY := &y
	var pointerZ *string
	fmt.Println(x, y, pointerX, pointerY, pointerZ)   // 10 true 0xc000020130 0xc000020134 <nil>
	fmt.Println(x, y, *pointerX, *pointerY, pointerZ) // 10 true 10 true <nil>
	// fmt.Println(*pointerZ) // panics, (you attempt to dereference a nil pointer)

	a := 10
	var pointerToA *int //  A pointer type can be based on any type.
	pointerToA = &a
	fmt.Println(pointerToA, a) // 0xc000020140 10

	var b = new(int)
	fmt.Println(b == nil) //  false
	fmt.Println(*b)       //  0

	// struct with a field of a pointer to a primitive type
	type person struct {
		FirstName  string
		MiddleName *string
		LastName   string
	}
	p := person{
		FirstName:  "Pat",
		MiddleName: stringp("Perry"), // This works, while:  "Perry"  or &"Perry"   will not
		LastName:   "Peterson",
	}
	fmt.Println(p)
}

func stringp(s string) *string {
	// Why does this work? When we pass a constant to a function, the constant is copied to a parameter, which is a variable.
	// Since it’s a variable, it has an address in memory. The function then returns the variable’s memory address.
	return &s
}
