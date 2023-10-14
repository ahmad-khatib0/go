package notifications

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/customers/customerspb"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/am"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/amotel"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/amprom"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/jetstream"
	pg "github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/postgres"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/postgresotel"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/registry"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/system"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/tm"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/notifications/internal/application"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/notifications/internal/constants"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/notifications/internal/grpc"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/notifications/internal/handlers"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/notifications/internal/postgres"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/ordering/orderingpb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = customerspb.Registrations(reg); err != nil {
		return err
	}
	if err = orderingpb.Registrations(reg); err != nil {
		return err
	}
	inboxStore := pg.NewInboxStore(constants.InboxTableName, svc.DB())
	messageSubscriber := am.NewMessageSubscriber(
		jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger()),
		amotel.OtelMessageContextExtractor(),
		amprom.ReceivedMessagesCounter(constants.ServiceName),
	)
	customers := postgres.NewCustomerCacheRepository(
		constants.CustomersCacheTableName,
		postgresotel.Trace(svc.DB()),
		grpc.NewCustomerRepository(svc.Config().Rpc.Service(constants.CustomersServiceName)),
	)

	// setup application
	app := application.New(customers)
	integrationEventHandlers := handlers.NewIntegrationEventHandlers(
		reg, app, customers,
		tm.InboxHandler(inboxStore),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, svc.RPC()); err != nil {
		return err
	}
	if err = handlers.RegisterIntegrationEventHandlers(messageSubscriber, integrationEventHandlers); err != nil {
		return err
	}

	return nil
}
