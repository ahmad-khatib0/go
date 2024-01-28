package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	order "github.com/ahmad-khatib0/go/grpc/resilient/8-tls-config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	order.UnimplementedOrderServiceServer
}

func getTlsCredentials() (credentials.TransportCredentials, error) {
	serverCert, serverCertErr := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if serverCertErr != nil {
		return nil, fmt.Errorf("could not load server key pairs: %s", serverCertErr)
	}

	// Since the CA signs the certificate, those shared certificates ( client.crt, server.crt) are already
	// in ca.crt. For development purposes, we will generate a cert pool in the server and client and append
	// client and server certificates there. Finally, we will put TLS configuration inside gRPC server options.
	certPool := x509.NewCertPool() // certificate pool for CA check

	caCert, caCertErr := os.ReadFile("cert/ca.crt")
	if caCertErr != nil {
		return nil, fmt.Errorf("could not read CA cert: %s", caCertErr)
	}

	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("failed to append the CA certs")
	}

	return credentials.NewTLS(
		&tls.Config{
			ClientAuth:   tls.RequireAnyClientCert,      // Client authentication type
			Certificates: []tls.Certificate{serverCert}, // provides server certificate
			ClientCAs:    certPool,                      // Roots the CA for the server to verify client certificates
		}), nil
}

func (s *server) Create(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	return &order.CreateOrderResponse{OrderId: 1243}, nil
}

func main() {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsCredentials, tlsCredentialsErr := getTlsCredentials()
	if tlsCredentialsErr != nil {
		log.Fatal("cannot load server TLS credentials: ", tlsCredentialsErr)
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(tlsCredentials))

	grpcServer := grpc.NewServer(opts...)
	order.RegisterOrderServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}
