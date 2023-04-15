package main

import "fmt"

const PublicMember = "public members starts with a capital latter"

func main() {
	var username string = "ahmad"
	fmt.Println(username)
	fmt.Printf("variable is of type %T \n", username)

	var isLoggedIn bool = true
	fmt.Println(isLoggedIn)
	fmt.Printf("variable is of type %T \n", isLoggedIn)

	var smallVal uint8 = 255
	fmt.Println(smallVal)
	fmt.Printf("variable is of type %T \n", smallVal)

	var smallFloat float64 = 255.4444444453
	fmt.Println(smallFloat)
	fmt.Printf("variable is of type %T \n", smallFloat)

	var zeroInitialization int
	fmt.Println(zeroInitialization) // 0
	fmt.Printf("variable is of type %T \n", zeroInitialization)

	var implicit = "auto detected type"
	fmt.Println(implicit)

	// without var keyword (valid only inside methods or functions, not globally)
	books := 330
	fmt.Println(books)

	fmt.Println(PublicMember)

}
