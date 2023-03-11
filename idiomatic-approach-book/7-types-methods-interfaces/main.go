package main

import "fmt"

func main() {
	// you can use any primitive type or compound type literal to define a concrete type
	primitiveTypes()

	// Method invocations
	p := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       33,
	}
	output := p.String()
	fmt.Println(output) // John Doe, age 33
}

func primitiveTypes() {
	type Score int
	type Converter func(string) Score
	type TeamScores map[string]Score
}

// Methods
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// Method declarations
func (p Person) String() string {
	return fmt.Sprintf("%s %s, age %d", p.FirstName, p.LastName, p.Age)
}
