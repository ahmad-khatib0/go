package main

import "fmt"

type Address struct {
	Street, City, Country string
}

type Person struct {
	Name    string
	Address *Address
}

func main() {
	john := Person{"John", &Address{"123 London Rd", "London", "UK"}}

	// jane := john                         //  shallow copy
	// jane.Name = "Jane"                   // ok
	// jane.Address.Street = "321 Baker St" //

	// fmt.Println(john.Name, john.Address) // John &{321 Baker St London UK}
	// fmt.Println(jane.Name, jane.Address) // Jane &{321 Baker St London UK}
	// so its not ok , because Address is a pointer

	// what you really want
	jane := john
	jane.Address = &Address{
		john.Address.Street,
		john.Address.City,
		john.Address.Country,
	}

	jane.Name = "Jane" // ok
	jane.Address.Street = "321 Baker St"

	fmt.Println(john.Name, john.Address)
	fmt.Println(jane.Name, jane.Address)
}
