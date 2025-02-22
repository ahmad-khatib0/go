package main

import (
	"context"
	"github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/resilient/2-ctx-timeout-propagation"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderServiceClient := order.NewOrderServiceClient(conn)
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)

	log.Println("Creating order...")
	_, errCreate := orderServiceClient.Create(ctx, &order.CreateOrderRequest{
		UserId:    23,
		ProductId: 123,
		Price:     12.3,
	})

	if errCreate != nil {
		log.Printf("Failed to create order. Err: %v", errCreate)
	} else {
		log.Println("Order is created successfully.")
	}
}
