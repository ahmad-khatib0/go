package main

import "fmt"

func main() {

	fmt.Println("structs in golang, there is no inheritance in golang , No super, no parent ...")

	ahmad := User{"ahmad", "ahmad@programmer.com", true, 19}
	fmt.Printf("ahmad details are: %+v\n ", ahmad) // %+v entire details ( key and value, )
	// while v% only for the value only
	fmt.Printf("name is %v , and email is %v", ahmad.Name, ahmad.Email)
}

// notice how everything is Capitalized letter to be public
type User struct {
	Name   string
	Email  string
	Status bool
	Age    int
}
