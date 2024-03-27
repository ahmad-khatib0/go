package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/baskets"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/cosec"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/customers"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/depot"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/config"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/system"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/web"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/migrations"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/notifications"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/ordering"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/payments"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/search"
	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/stores"
)

type monolith struct {
	*system.System
	modules []system.Module
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("mallbots exitted abnormally: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	var cfg config.AppConfig
	cfg, err = config.InitConfig()
	if err != nil {
		return err
	}

	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}

	m := monolith{
		System: s,
		modules: []system.Module{
			&baskets.Module{},
			&customers.Module{},
			&depot.Module{},
			&notifications.Module{},
			&ordering.Module{},
			&payments.Module{},
			&stores.Module{},
			&cosec.Module{},
			&search.Module{},
		},
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(m.DB())

	err = m.MigrateDB(migrations.FS)
	if err != nil {
		return err
	}

	if err = m.startupModules(); err != nil {
		return err
	}

	// Mount general web resources
	m.Mux().Mount("/", http.FileServer(http.FS(web.WebUI)))

	fmt.Println("started mallbots application")
	defer fmt.Println("stopped mallbots application")

	m.Waiter().Add(m.WaitForWeb, m.WaitForRPC, m.WaitForStream)

	// go func() {
	// 	for {
	// 		var mem runtime.MemStats
	// 		runtime.ReadMemStats(&mem)
	// 		m.logger.Debug().Msgf("Alloc = %v  TotalAlloc = %v  Sys = %v  NumGC = %v", mem.Alloc/1024, mem.TotalAlloc/1024, mem.Sys/1024, mem.NumGC)
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	return m.Waiter().Wait()
}

func (m *monolith) startupModules() error {
	for _, module := range m.modules {
		ctx := m.Waiter().Context()
		if err := module.Startup(ctx, m); err != nil {
			return err
		}
	}

	return nil
}
