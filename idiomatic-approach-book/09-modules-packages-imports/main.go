package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
)

func main() {
	seedRand()
}

func seedRand() *rand.Rand {
	var b [8]byte
	fmt.Println(b)
	_, err := crand.Read(b[:])
	if err != nil {
		panic("cannot seed with cryptographic random number generator")
	}
	r :=
		rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
	return r
}

// ***************************  Gracefully Renaming and Reorganizing Your API  *************************

type Foo struct {
	x int
	S string
}

func (f Foo) Hello() string {
	return "hello"
}
func (f Foo) goodbye() string {
	return "goodbye"
}

type Bar = Foo //  to access Foo by the name Bar

func MakeBar() Bar {
	bar := Bar{
		x: 20,
		S: "Hello",
	}
	var f Foo = bar
	fmt.Println(f.Hello())
	return bar
}

func renamingTypes() {
	type Bar = Foo //  to access Foo by the name Bar

}
