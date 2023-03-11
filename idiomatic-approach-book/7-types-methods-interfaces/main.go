package main

import (
	"fmt"
	"time"
)

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

	var c Counter
	fmt.Println(c.String()) // total: 0, last updated: 0001-01-01 00:00:00 +0000 UTC
	c.Increment()
	fmt.Println(c.String()) // total: 1, last updated: 2009-11-10 23:00:00 +0000 UTC m=+0.000000001

	methodsForNilInstances()
}

func primitiveTypes() {
	type Score int
	type Converter func(string) Score
	type TeamScores map[string]Score
}

// *******************************  Methods   ****************************************
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// Method declarations
func (p Person) String() string {
	return fmt.Sprintf("%s %s, age %d", p.FirstName, p.LastName, p.Age)
}

// Pointer Receivers and Value Receivers
type Counter struct {
	total       int
	lastUpdated time.Time
}

func (c *Counter) Increment() {
	c.total++
	c.lastUpdated = time.Now()
}
func (c *Counter) String() string {
	return fmt.Sprintf("total: %d last updated %v", c.total, c.lastUpdated)
}

// Code Your Methods for nil Instances
type IntTree struct {
	val         int
	left, right *IntTree
}

func (it *IntTree) Insert(val int) *IntTree {
	if it == nil {
		return &IntTree{val: val}
	}
	if val < it.val {
		it.left = it.left.Insert(val)
	} else if val > it.val {
		it.right = it.right.Insert(val)
	}
	return it
}

func (it *IntTree) Contains(val int) bool {
	switch {
	case it == nil:
		return false
	case val < it.val:
		return it.left.Contains(val)
	case val > it.val:
		return it.right.Contains(val)
	default:
		return true
	}
}

// The Contains method doesn’t modify the *IntTree, but it is declared with a pointer receiver. This demonstrates
// the rule mentioned above about supporting a nil receiver. A method with a value receiver ( can’t check for nil )
// and as mentioned earlier, panics if invoked with a nil receiver.
func methodsForNilInstances() {
	var it *IntTree
	it = it.Insert(5)
	it = it.Insert(3)
	it = it.Insert(10)
	it = it.Insert(2)
	fmt.Println(it.Contains(2))  // true
	fmt.Println(it.Contains(12)) // false

}
