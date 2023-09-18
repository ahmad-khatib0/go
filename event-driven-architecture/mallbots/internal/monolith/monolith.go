package monolith

import (
	"context"
	"database/sql"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/config"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/waiter"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Monolith interface {
	DB() *sql.DB
	Logger() zerolog.Logger
	Mux() *chi.Mux
	RPC() *grpc.Server
	Config() config.AppConfig
	Waiter() waiter.Waiter
}

type Module interface {
	Startup(context.Context, Monolith) error
}
