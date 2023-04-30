package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// +------------------------------------------------------------------------------------------------+
// | with mutex , once we use the go keyword, we can not to talk to that routine , we can only wait |
// | for it to finishو, but with channels we can to talk to that routine, and exchange data with it |
// +------------------------------------------------------------------------------------------------+

const NumberOfPizzas = 10

var PizzasMade, PizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error // one channel that holds a channel of errors
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch // will return nil if channel closed successfully , and error if any
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making

	// run forever or until we receive a quit notification
	// try to make pizzas
	for {
		// try to make a pizza
		// decision
	}

}

func main() {

	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	color.Cyan("The Pizzeria is open for business!")
	color.Cyan("----------------------------------")

	// create a producer
	pizzaJop := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria(pizzaJop)

	// create and run consumer

	// print out the ending message
}
