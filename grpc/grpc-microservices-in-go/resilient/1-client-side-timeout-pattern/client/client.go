package main

import (
	"context"
	shipping "github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/resilient/1-client-side-timeout-pattern"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect shipping service. Err: %v", err)
	}

	defer conn.Close()

	shippingClient := shipping.NewShippingServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	log.Println("Creating shipping...")

	_, errCreate := shippingClient.Create(ctx, &shipping.CreateShippingRequest{UserId: 23})
	if errCreate != nil {
		log.Printf("Failed to create shipping. Err: %v", errCreate)
	} else {
		log.Println("Shipping is created successfully.")
	}
}
