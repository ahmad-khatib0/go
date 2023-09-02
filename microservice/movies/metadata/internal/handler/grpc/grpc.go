package grpc

import (
	"context"

	"github.com/ahmad-khatib0/go/microservice/movies/gen"
	"github.com/ahmad-khatib0/go/microservice/movies/metadata/internal/controller/metadata"
	"github.com/ahmad-khatib0/go/microservice/movies/metadata/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a movie metadata gRPC handler.
type Handler struct {
	gen.UnimplementedMetadataServiceServer // This is required by a Protocol Buffers compiler to enforce future compatibility.
	ctrl                                   *metadata.Controller
}

// New creates a new movie metadata gRPC handler.
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMetadataByID returns movie metadata by id.
func (h *Handler) GetMetadataByID(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id")
	}

	m, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}
