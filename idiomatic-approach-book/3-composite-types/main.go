package main

import "fmt"

func arrays() {
	var a [3]int // [0.0.0]
	b := [3]int{10, 20, 30}

	// If you have a sparse array (an array where most elements are set to their zero value),
	// you can specify only the indices with values in the array literal:
	c := [12]int{1, 5: 4, 6, 10: 100, 15} // [1, 0, 0, 0, 0, 4, 6, 0, 0, 0, 100, 15]

	// When using an array literal to initialize an array, you can leave off the number:
	d := [...]int{10, 20, 30}
	fmt.Println(b == d) // true

	// there are languages with true matrix support; Go isn’t one of them
	// Go only has one-dimensional arrays, but you can simulate multi-dimensional arrays:
	var e [2][3]int // [[0 0 0] [0 0 0]]

	fmt.Println(a, b, c, e)

	fmt.Println(len(e)) // 2
}

func slices() {
	a := []int{10, 20, 30}
	b := []int{1, 5: 4, 6, 10: 100, 15} //  [1, 0, 0, 0, 0, 4, 6, 0, 0, 0, 100, 15]

	var c [][]int //  multi-dimensional slices  => []
	var d []int   //  A nil slice contains nothing => []

	// slices are comparable with nil only
	fmt.Println(c == nil) // true

	d = append(d, 4, 5, 6)
	d = append(d, a...)

	fmt.Println(a, b, c, d)
}

func capacity() {
	// The built-in cap function returns the current capacity of a slice. It is used far less frequently than len.
	// Most of the time, cap is used to check if a slice is large enough to
	// hold new data, or if a call to make is needed to create a new slice
	var a []int
	fmt.Println(a, len(a), cap(a)) // [] 0 0
	a = append(a, 10)
	fmt.Println(a, len(a), cap(a)) // [10] 1 1
	a = append(a, 20)
	fmt.Println(a, len(a), cap(a)) // [10 20] 2 2
	a = append(a, 30)
	fmt.Println(a, len(a), cap(a)) // [10 20 30] 3 4
	a = append(a, 40)
	fmt.Println(a, len(a), cap(a)) // [10 20 30 40] 4 4
	a = append(a, 50)
	fmt.Println(a, len(a), cap(a)) // [10 20 30 40 50] 5 8   => it's doubling the capacity

	// so While it’s nice that slices grow automatically, it’s far more efficient to
	// size them once. If you know how many things you plan to put into a slice, create the
	// slice with the correct initial capacity. We do that with the make function.
}

func makesAndSlices() {
	// This creates an int slice with a length of 5 and a capacity of 5
	x := make([]int, 5) // [0 0 0 0 0]

	x = append(x, 10) // [0 0 0 0 0 10]  length of 6 and a capacity of 10
	// BECAUSE APPEND ALWAYS INCREASES THE LENGTH OF A SLICE

	a := make([]int, 0, 10) // we have a non-nil slice with a length of zero, but a capacity of 10
	a = append(a, 5, 6, 7, 8)

	// Slicing slices
	b := []int{1, 2, 3, 4}
	c := b[:2] // leave off the starting offset, 0 is assumed.
	d := b[1:] // leave off the ending offset, the end of the slice is substituted
	e := b[1:3]
	f := b[:]
	fmt.Println("b:", b) // b: [1 2 3 4]
	fmt.Println("c:", c) // c: [1 2]
	fmt.Println("d:", d) // d: [2 3 4]
	fmt.Println("e:", e) // e: [2 3]
	fmt.Println("f:", f) // f: [1 2 3 4]

	// NOTE: changes to an element in a slice affect all slices that share that element
	g := []int{1, 2, 3, 4}
	h := g[:2]
	j := g[1:]
	g[1] = 20
	h[0] = 10
	j[1] = 30
	fmt.Println("g:", g) // g: [10 20 30 4]
	fmt.Println("h:", h) // h: [10 20]
	fmt.Println("j:", j) // j: [20 30 4]

	k := []int{1, 2, 3, 4}
	l := k[:2]
	fmt.Println(cap(k), cap(l)) // 4 4    => shared capacity
	l = append(l, 30)
	fmt.Println("k:", k) // k: [1 2 30 4]  => changed 3th element
	fmt.Println("l:", l) // y: [1 2 30]

	// even more confusing
	m := make([]int, 0, 5)
	m = append(m, 1, 2, 3, 4)
	n := m[:2]
	o := m[2:]
	fmt.Println(cap(m), cap(n), cap(o)) // 5 5 3
	n = append(n, 30, 40, 50)
	m = append(m, 60)
	fmt.Println("n:", n) //  [1 2 30 40 60]
	fmt.Println("m:", m) //  [1 2 30 40 60]
	o = append(o, 70)
	fmt.Println("m:", m) //  [1 2 30 40 70]
	fmt.Println("n:", n) //  [1 2 30 40 70]
	fmt.Println("o:", o) //  [30 40 70]
}

func fullSliceExpression() {
	// The full slice expression protects against append overwritting (as the last tow examples in makesAndSlices )
	x := make([]int, 0, 5)
	x = append(x, 1, 2, 3, 4)
	y := x[:2:2]
	z := x[2:4:4]
	fmt.Println(cap(x), cap(y), cap(z)) // 5 2 2
	y = append(y, 30, 40, 50)
	x = append(x, 60)
	z = append(z, 70)
	fmt.Println("x:", x) // x: [1 2 3 4 60]
	fmt.Println("y:", y) // y: [1 2 30 40 50]
	fmt.Println("z:", z) // z: [3 4 70]

	// Because we limited the capacity of the subslices to their lengths, appending additional elements
	// onto y and z created new slices that didn’t interact with the other slices. After this code runs
}

func convertArraysToSlices() {
	//  that taking a slice from an array has the same memory-sharing properties as taking a slice from a slice
	x := [4]int{5, 6, 7, 8}
	y := x[:2]
	z := x[2:]
	x[0] = 10
	fmt.Println("x:", x) // x: [10 6 7 8]
	fmt.Println("y:", y) // y: [10 6]   => slice overwritting
	fmt.Println("z:", z) // z: [7 8]
}

func copySlice() {
	// If you need to create a slice that’s independent of the original
	x := []int{1, 2, 3, 4}
	y := make([]int, 4)
	num := copy(y, x)   // destination , source
	fmt.Println(y, num) // [1 2 3 4] 4

	y = make([]int, 2)
	num = copy(y, x)    // copies the first two elements of a four element slice
	fmt.Println(y, num) // [1 2 ] 2

	copy(y, x[2:]) // copy from the middle of the source slice
	// Also note that we don’t assign the output of copy to a variable. If you don’t
	// need the number of elements copied, you don’t need to assign it

	// The copy function allows you to copy between two slices that cover overlapping sections of an underlying slice:
	x = []int{1, 2, 3, 4}
	num = copy(x[:3], x[1:])
	fmt.Println(x, num) // [2 3 4 4] 3

	// You can use copy with arrays by taking a slice of the array.
	// You can make the array either the source or the destination of the copy
	x = []int{1, 2, 3, 4}
	d := [4]int{5, 6, 7, 8}
	y = make([]int, 2)
	copy(y, d[:])
	fmt.Println(y) // [5 6]
	copy(d[:], x)
	fmt.Println(d) // [1 2 3 4]   => overwrote all the d elements with the x ones
}

func slicingStrings() {
	var s string = "Hello there"
	var s2 string = s[4:7]
	var s3 string = s[:5]
	var s4 string = s[6:]
	fmt.Println(s2, s3, s4) // o t Hello there

	s = "Hello "
	s2 = s[4:7]
	s3 = s[:5] // s3 is not set to “o ”. Instead, you get “o �” That’s because
	// we only copied the first byte of the sun emoji’s code point
	s4 = s[6:]

	s = "Hello "
	fmt.Println(len(s))
	// This code prints out 10, not 7, because it takes 4 bytes to represent the sun with smiling face emoji in UTF-8.

	var a rune = 'x'
	s = string(a) // A single rune or byte can be converted to a string:
	var b byte = 'y'
	s2 = string(b)

	s = "Hello, "
	var bs []byte = []byte(s)
	var rs []rune = []rune(s)
	fmt.Println(bs, rs) // [72 101 108 108 111 44 32 240 159 140 158] [72 101 108 108 111 44 32 127774]

}

func main() {
	// arrays()
	// slices()
	// capacity()
	makesAndSlices()
	copySlice()
	slicingStrings()
}
