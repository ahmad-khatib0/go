package main

import (
	"context"
	order "github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/resilient/5-client-circuit-breaker"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"log"
	"time"
)

var cb *gobreaker.CircuitBreaker

// this logic to a client-server example in which the client tries to create an order periodically, and an
// exception for some requests is thrown due to an outage in the Order service. For this
// kind of error case, applying the circuit breaker pattern to client-server interaction is
// beneficial instead of blindly sending requests
func main() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "demo",
		MaxRequests: 3,
		Timeout:     4,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return failureRatio >= 0.1
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit Breaker: %s, changed from %v, to %v", name, from, to)
		},
	})

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderServiceClient(conn)
	for {
		log.Println("Creating order...")

		// rapping the order service call with circuit breaker
		orderResponse, errCreate := cb.Execute(func() (interface{}, error) {
			return orderClient.Create(context.Background(), &order.CreateOrderRequest{
				UserId:    23,
				ProductId: 1234,
				Price:     3.2,
			})
		})

		if errCreate != nil {
			// Error when state changes
			log.Printf("Failed to create order. Err: %v", errCreate)
		} else {
			// If the circuit is closed, it returns data
			log.Printf("Order %d is created successfully.", orderResponse.(*order.CreateOrderResponse).OrderId)
		}

		time.Sleep(1 * time.Second)
	}
}
