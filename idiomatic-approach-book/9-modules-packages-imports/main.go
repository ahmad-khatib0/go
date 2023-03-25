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
