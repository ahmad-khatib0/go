package main

import (
	"context"
	order "github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/resilient/7-server-client-validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderServiceClient(conn)
	log.Println("Creating order...")
	orderResponse, errCreate := orderClient.Create(context.Background(), &order.CreateOrderRequest{
		UserId:    -1,
		ProductId: 0,
		Price:     2,
	})

	if errCreate != nil {
		stat := status.Convert(errCreate)
		for _, detail := range stat.Details() {
			switch errType := detail.(type) {
			case *errdetails.BadRequest:
				for _, violation := range errType.GetFieldViolations() {
					log.Printf("The field %s has invalid value. desc: %v", violation.GetField(), violation.GetDescription())
				}
			}
		}

	} else {
		log.Printf("Order %d is created successfully.", orderResponse.OrderId)
	}
}
