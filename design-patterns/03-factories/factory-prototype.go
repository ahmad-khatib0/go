package main

import "fmt"

type Employee struct {
	Name, Position string
	AnnualIncome   int
}

// create multiple constants representing predefined Employees prototypes (or templates )
const (
	Developer = iota
	Manager
)

func NewEmployee(role int) *Employee {
	switch role {
	case Developer:
		return &Employee{"", "developer", 50000}
	case Manager:
		return &Employee{"", "manager", 80000}
	default:
		panic("unsupporter role")
	}
}

func main() {
	dev := NewEmployee(Developer)
	dev.Name = "Ahmad"

	fmt.Println(dev)
}
