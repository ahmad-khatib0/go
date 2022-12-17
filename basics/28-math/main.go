package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"time"
)

func main() {
	fmt.Println("math and crypto in golang")

	// var firstNum  int = 3
	// var secondNum  float64 = 4.5
	// fmt.Println("the sum of tow numbers is: " , firstNum + int(secondNum))  // sure incorrect

	// random numbers
	mathRand.Seed(time.Now().UnixNano()) // without this random seed, rand.Intn will print alwasy the same resutl
	fmt.Println(mathRand.Intn(5) + 1)

	// randomness from crypto
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(5)) // between 0 - 5
	fmt.Println(randomNum)
}
