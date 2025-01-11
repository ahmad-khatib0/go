package main

import "log"

type Person struct {
	Name     string
	Age      int
	EeyCount int
}

// sometimes you want additional logic, additional stuff to happen as you are
// creating a particular struct. And that's when you would use these factory functions.

func NewPerson(name string, age int) *Person {
	if age < 16 {

	}
	return &Person{name, age, 3}
}

func main() {
	p := NewPerson("test", 33)
	log.Println(p)
}
