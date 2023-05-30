package main

import "fmt"

func main() {
	fmt.Println("pointers ")

	// var pointer *int   // creating a pointer
	// fmt.Print("value of pointer is: ", pointer)

	myPointer := 33
	var ptr = &myPointer // create a pointer with a reference to other variables

	fmt.Println("value of actual pointer is: ", ptr)  //  0xc0000200f0 the address in memory
	fmt.Println("value of actual pointer is: ", *ptr) // 33 , invoking that created pointer's value

	*ptr = *ptr + 2
	fmt.Println("new value is: ", myPointer) // 35 because myPointer references that ptr
}
