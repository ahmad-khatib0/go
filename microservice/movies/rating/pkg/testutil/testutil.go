package testutil

import (
	"github.com/ahmad-khatib0/go/microservice/movies/gen"
	"github.com/ahmad-khatib0/go/microservice/movies/rating/internal/controller/rating"
	grpchandler "github.com/ahmad-khatib0/go/microservice/movies/rating/internal/handler/grpc"
	"github.com/ahmad-khatib0/go/microservice/movies/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server to be used in tests.
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}
