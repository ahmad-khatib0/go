package main

import (
	"context"
	"log"

	payment "github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/basic-server-client/proto/payment"

	"google.golang.org/grpc"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("http:/localhost:8080", opts...)
	if err != nil {
		log.Println("It is fine, this is not a complete example.")
	}

	defer conn.Close()

	paymentClient := payment.NewPaymentServiceClient(conn)
	ctx := context.Background()
	_, err = paymentClient.Create(ctx, &payment.CreatePaymentRequest{Price: 23})
	if err != nil {
		log.Println("Don't worry, we don't expect to see it is working.")
	}
}
