package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func main() {
	fmt.Println("interfaces")
	c := Client{L: LogicProvider{}}
	c.Progeram()

	theAnyType()
	typeAssertion()
}

type LogicProvider struct{}

func (lp LogicProvider) Process(data string) string {
	return data
}

type Logic interface {
	Process(data string) string
}

type Client struct {
	L Logic
}

func (c Client) Progeram() {
	fmt.Println(c.L.Process("data passed to process"))
}

// In the Go code, there is an interface, but only the caller (Client) knows about it; there is nothing
// declared on LogicProvider to indicate that it meets the interface. This is sufficient to both allow a
// new logic provider in the future as well as provide executable documentation to ensure that any type
// passed in to the client will match the client’s need.

//  ╭──────────────────────────────────────────────────────────────────╮
//  │      Interfaces specify what callers need. The client code       │
//  │ defines the interface to specify what functionality it requires. │
//  ╰──────────────────────────────────────────────────────────────────╯

// Embedding and Interfaces

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

type ReadCloser interface {
	Reader
	Closer
}

// Interfaces and nil
func interfacesAndNil() {
	// In order for an interface to be considered nil both the type and the value must be nil
	var s *string
	fmt.Println(s == nil) // prints true
	var i interface{}
	fmt.Println(i == nil) // prints true
	i = s
	fmt.Println(i == nil) // prints false
}

// the Any type by using an (Empty Interface)
func theAnyType() {
	var i interface{}
	i = 20
	i = "hello"
	i = struct {
		FirstName string
		LastName  string
	}{"Fred", "Fredson"}
	fmt.Println(i) // {Fred Fredson}
}

// using interface  is as way to store a value in a user-created data structure.
// This is due to Go’s current lack of user-defined generics
func fileReader(filename string) error {
	data := map[string]interface{}{}
	contents, err := ioutil.ReadFile("testdata/sample.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(contents, &data)
}

type LinkedList struct {
	Value interface{}
	Next  *LinkedList
}

func (ll *LinkedList) Insert(pos int, val interface{}) *LinkedList {
	if ll == nil || pos == 0 {
		return &LinkedList{
			Value: val,
			Next:  ll,
		}
	}

	ll.Next = ll.Next.Insert(pos-1, val)
	return ll
}

// Type Assertions and Type Switches
type MyInt int

func typeAssertion() (int, error) {
	var i interface{}
	var mine MyInt = 20
	i = mine

	i2, ok := i.(int)
	if !ok {
		return 0, fmt.Errorf("unexpected type for %v", i)
	}
	return i2 + 1, nil
}

func typeSwitches(i interface{}) {
	switch j := i.(type) {
	case nil: // i is nil, type of j is interface{}
	case int: // j is of type int
	case MyInt: // j is of type MyInt
	case io.Reader: // j is of type io.Reader
	case string: // j is a string
	case bool, rune: // i is either a bool or rune, so j is of type interface{}
	default:
		fmt.Println(j) // // no idea what i is, so j is of type interface{}
	}

}
