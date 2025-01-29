package main

import (
	"errors"
	"github.com/sony/gobreaker"
	"log"
	"math/rand"
)

var cb *gobreaker.CircuitBreaker

func main() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3, // Allowed number of requests for a half-open circuit
		Timeout:     4, // timeout for an open to half-open transition
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Decides on if the circuit will be open
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.6
		},

		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})

	cbRes, cbErr := cb.Execute(func() (interface{}, error) {
		// Wrapped function to apply circuit breaker
		res, isErr := isError()
		if isErr {
			return nil, errors.New("error")
		}
		return res, nil
	})

	if cbErr != nil {
		// Returns an error once the circuit is open
		log.Fatalf("Circuit breaker error %v", cbErr)
	} else {
		log.Printf("Circuit breaker result %v", cbRes)
	}
}

func isError() (int, bool) {
	min := 10
	max := 30
	result := rand.Intn(max-min) + min
	return result, result != 20
}
