package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	welcome := "user inputs"
	fmt.Println(welcome)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("please provide a rating for our product")

	// comma Ok  || error Ok syntax
	input, _ := reader.ReadString('\n')
	fmt.Println("Thanks for your rating: ", input)

}
