package main

import "fmt"

func main() {

	ahmad := User{"ahmad", "ahmad@programmer.com", true, 19}
	fmt.Printf("ahmad details are: %+v\n ", ahmad)
	fmt.Printf("name is %v , and email is %v", ahmad.Name, ahmad.Email) // ahmad@programmer.com

	ahmad.GetStatus()
	ahmad.AssignNewEmail()                                              // new@test.com
	fmt.Printf("name is %v , and email is %v", ahmad.Name, ahmad.Email) // ahmad@programmer.com

	// NOTE:  it DID NOT change the actual first initialized value
}

// notice how everything is Capitalized letter to be public
type User struct {
	Name   string
	Email  string
	Status bool
	Age    int
}

// define a method
func (u User) GetStatus() {
	fmt.Println("is user active? ", u.Status)
}

func (u User) AssignNewEmail() {
	u.Email = "new@test.com"
	fmt.Println("the new email is: ", u.Email)
}
