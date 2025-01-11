package main

import (
	"fmt"
)

type Person interface {
	SayHello()
}

type person struct {
	name string
	age  int
}

type tiredPerson struct {
	name string
	age  int
}

// ╒═══════════════════════════════════════════════════════════════════════════════════╕
//
//	so with factory interface we can have multiple types , so the actual
//	underlying object is different depending on what invocation data you use here
//	AND ALSO BECAUSE STRUCTS FIELDS ARE PRIVATE, SO THEY ARE PROTECTED FROM MODIFYING
//
// └───────────────────────────────────────────────────────────────────────────────────┘

func NewPerson(name string, age int) Person {
	if age > 60 {
		return &tiredPerson{name, age}
	}
	return &person{name, age}
}

func (p *person) SayHello() {
	fmt.Printf("Hello my name is %s , and my age is %d ", p.name, p.age)
}

func (p *tiredPerson) SayHello() {
	fmt.Println("sorry i'm too tired to talk to you")
}

// So this is a neat way of encapsulating some information and just having your
// factory expose just the interface that you can subsequently work with.
// And that way you can you can, for example, have different underlying types

func main() {
	p := NewPerson("june", 33)
	p.SayHello()
}
