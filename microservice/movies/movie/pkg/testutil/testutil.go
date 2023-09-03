package testutil

import (
	"github.com/ahmad-khatib0/go/microservice/movies/gen"
	"github.com/ahmad-khatib0/go/microservice/movies/movie/internal/controller/movie"
	metadatagateway "github.com/ahmad-khatib0/go/microservice/movies/movie/internal/gateway/metadata/grpc"
	ratinggateway "github.com/ahmad-khatib0/go/microservice/movies/movie/internal/gateway/rating/grpc"
	grpchandler "github.com/ahmad-khatib0/go/microservice/movies/movie/internal/handler/grpc"
	"github.com/ahmad-khatib0/go/microservice/movies/pkg/discovery"
)

// NewTestMovieGRPCServer creates a new movie gRPC server to be used in tests.
func NewTestMovieGRPCServer(registry discovery.Registry) gen.MovieServiceServer {
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	return grpchandler.New(ctrl)
}
