package baskets

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/application"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/grpc"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/handlers"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/logging"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/postgres"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets/internal/rest"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/ddd"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/monolith"
)

type Module struct{}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) (err error) {
	// setup Driven adapters
	domainDispatcher := ddd.NewEventDispatcher()
	baskets := postgres.NewBasketRepository("baskets.baskets", mono.DB())
	conn, err := grpc.Dial(ctx, mono.Config().Rpc.Address())
	if err != nil {
		return err
	}
	stores := grpc.NewStoreRepository(conn)
	products := grpc.NewProductRepository(conn)
	orders := grpc.NewOrderRepository(conn)

	// setup application
	app := logging.LogApplicationAccess(
		application.New(baskets, stores, products, orders, domainDispatcher),
		mono.Logger(),
	)
	orderHandlers := logging.LogDomainEventHandlerAccess(
		application.NewOrderHandlers(orders),
		mono.Logger(),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(app, mono.RPC()); err != nil {
		return err
	}

	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}

	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	handlers.RegisterOrderHandlers(orderHandlers, domainDispatcher)

	return
}
