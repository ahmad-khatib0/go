package main

import (
	"fmt"
	"log"
	"net"

	payment "github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/basic-server-client/proto/payment"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	payment.RegisterPaymentServiceServer(grpcServer, payment.UnimplementedPaymentServiceServer{})
	grpcServer.Serve(listener)
}
