package main

// +------------------------------------------------------------------------------------------------+
// | with mutex , once we use the go keyword, we can not to talk to that routine , we can only wait |
// | for it to finishو, but with channels we can to talk to that routine, and exchange data with it |
// +------------------------------------------------------------------------------------------------+

const NumberOfPizzas = 10

var PizzasMade, PizzasFailed, total int

type Porducer struct {
	data chan PizzaOrder
	quit chan chan error // one channel that holds a channel of errors
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func main() {

	// seed the random number generator

	// print out a message

	// create a producer

	// run the producer in the background

	// create and run consumer

	// print out the ending message
}
