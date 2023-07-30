package main

import (
	"fmt"
	"log"
)

type Employee struct {
	Name, Position string
	AnnualIncome   int
}

// we want to be able to create factories dependent upon the settings
// that we want employees to subsequently be manufactured.

// FUNCTIONAL
// So notice we're not creating an employee. We're creating an employee factory that you
// can subsequently use to fine tune those details of that object.
func NewEmployeeFactory(position string, annualIncome int) func(name string) *Employee {
	return func(name string) *Employee {
		return &Employee{name, position, annualIncome}
	}
}

// /
// /
// /
// /
// The functional approach is pretty good, but there is another approach, of course,
// and that is a more structural approach, basically making a factory a struct.
type EmployeeFactory struct {
	Position     string
	AnnualIncome int
}

// make a factory function for actually returning instances of this particular factory
// and you would also give this factory some sort of method for actually creating the object.
func (e *EmployeeFactory) Create(name string) *Employee {
	return &Employee{name, e.Position, e.AnnualIncome}
}

// what you need to do is you need to have a factory function for creating this factory because.
// Well, because you want to have different predefined employee factories, obviously.
func NewEmployeeFactory2(position string, annualIncome int) *EmployeeFactory {
	return &EmployeeFactory{position, annualIncome}
}

func main() {
	//  FUNCTIONAL
	//  ╒════════════════════════════════════════════════════════════════════════════════╕
	//    one advantage is that now you have these factories stored in variables,
	//    you can also pass these variables into other functions. And that is, you know,
	//    the core of functional programming, passing functions into functions.
	//  └────────────────────────────────────────────────────────────────────────────────┘

	developerFactory := NewEmployeeFactory("developer", 4000)
	managerFactory := NewEmployeeFactory("manager", 6000)

	developer := developerFactory("Anna")
	manager := managerFactory("Tom")

	log.Println(developer)
	log.Println(manager)

	//
	bossFactory := NewEmployeeFactory2("CEO", 100000)
	boss := bossFactory.Create("Jacoup")
	fmt.Println(boss)
	// the addetional advantage over the functional factory is that here you
	// can modify the created factories later on, for example:
	boss.AnnualIncome = 110000
}
