package main

import (
	"context"

	pb "github.com/ahmad-khatib0/go/grpc/grpc-up-and-running/ch05/interceptors/order-service/go/order-service-gen"
	"google.golang.org/grpc/encoding/gzip"

	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Setting up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// RPC: Add Order
	order1 := pb.Order{Id: "101", Items: []string{"iPhone XS", "Mac Book Pro"}, Destination: "San Jose, CA", Price: 2300.00}
	res, _ := client.AddOrder(ctx, &order1, grpc.UseCompressor(gzip.Name))

	log.Print("AddOrder Response -> ", res.Value)

}
