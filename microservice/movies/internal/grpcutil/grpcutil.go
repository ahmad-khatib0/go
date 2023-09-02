package grpcutil

import (
	"context"
	"math/rand"

	"github.com/ahmad-khatib0/go/microservice/movies/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceConnection attempts to select a random service
// instance and returns a gRPC connection to it.
func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry) (*grpc.ClientConn, error) {

	adrs, err := registry.ServiceAddresses(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(adrs[rand.Intn(len(adrs))], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
