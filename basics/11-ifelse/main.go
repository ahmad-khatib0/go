package main

import "fmt"

func main() {

	userAge := 33
	var result string

	if userAge > 15 {
		result = "regular user"
	} else if userAge == 18 {
		result = "also allowed"
	} else {
		result = "under age"
	}

	fmt.Println("user status: ", result)

	// assign a variable and make a condition
	if num := 4; num < 10 {
		fmt.Println("number is less than 4")
	} else {
		fmt.Println("number is greater than 4")
	}
}
