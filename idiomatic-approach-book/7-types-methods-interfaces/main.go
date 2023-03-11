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

	methodsAreFunctions()
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

// Methods Are Functions Too
type Adder struct {
	start int
}

func (a Adder) AddTo(val int) int {
	return a.start + val
}

func methodsAreFunctions() {
	myAdder := Adder{start: 10}
	fmt.Println(myAdder.AddTo(5)) // prints 15

	f1 := myAdder.AddTo // We can also assign the method to a variable or pass it to
	// a parameter of type func(int)int. This is called a method value:
	fmt.Println(f1(10)) // prints 20

	f2 := Adder.AddTo            // You can also create a function from the type itself. This is called a method expression
	fmt.Println(f2(myAdder, 15)) // prints 25
}

func iotaIsForEnumerations() {
	type MailCategory int
	const (
		UnCategorized MailCategory = iota
		Personal
		Spam
		Social
		Advertisements
	)
	// The first constant in the const block has the type specified and its value is set to iota. Every subsequent
	// line has neither the type nor a value assigned to it. When the Go compiler sees this, it repeats the type
	// and the assignment to all of the subsequent constants in the block, and increments the value of iota on each
	// line. This means that it assigns 0 to the first constant (UnCategorized), 1 to the second constant (Personal)
	// and so on. When a new const block is created, iota is set back to 0
}
