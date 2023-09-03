package testutil

import (
	"github.com/ahmad-khatib0/go/microservice/movies/gen"
	"github.com/ahmad-khatib0/go/microservice/movies/metadata/internal/controller/metadata"
	grpchandler "github.com/ahmad-khatib0/go/microservice/movies/metadata/internal/handler/grpc"
	"github.com/ahmad-khatib0/go/microservice/movies/metadata/internal/repository/memory"
)

// NewTestMetadataGRPCServer creates a new metadata gRPC server to be used in tests.
func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	r := memory.New()
	ctrl := metadata.New(r)
	return grpchandler.New(ctrl)
}
