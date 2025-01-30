package gapi

import (
	"fmt"

	db "github.com/ahmad-khatib0/go/simple-bank/project/db/sqlc"
	"github.com/ahmad-khatib0/go/simple-bank/project/pb"
	"github.com/ahmad-khatib0/go/simple-bank/project/token"
	"github.com/ahmad-khatib0/go/simple-bank/project/util"
	"github.com/ahmad-khatib0/go/simple-bank/project/worker"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	// pb.UnimplementedSimpleBankServerIts main purpose is to enable forward compatibility, Which means that the
	// server can already accept the calls to the CreateUser and LoginUser RPCs before they are
	// actually implemented. Then we can gradually add their real implementations later.
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
