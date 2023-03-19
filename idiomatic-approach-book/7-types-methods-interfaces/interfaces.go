package main

import "fmt"

func main() {
	fmt.Println("interfaces")
	c := Client{L: LogicProvider{}}
	c.Progeram()
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
