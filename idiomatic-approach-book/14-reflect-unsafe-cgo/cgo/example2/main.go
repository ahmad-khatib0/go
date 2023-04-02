package main

import "fmt"

/*
	extern int add(int a, int b);
*/
import "C"

func main() {
	sum := C.add(3, 2)
	fmt.Println(sum)
}

//  Go function can be exposed to C code by putting an //export comment before the function:

//export doubler
func doubler(i int) int {
	return i * 2
}
