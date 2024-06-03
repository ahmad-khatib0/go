package main

import (
	"log"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/profile"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/push"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/server"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"go.uber.org/zap/zapcore"
)

type app struct {
	logger    *logger.Logger
	cfg       *config.Config
	store     *store.Store
	statsChan *stats.Stats
	utils     *utils.Utils
	cluster   *server.Cluster
	profile   *profile.Profile

	authValidators  map[types.Level][]string
	validators      map[string]server.CredValidator
	validatorCliCfg map[string][]string
	// Tag namespaces (prefixes) which are immutable to the client.
	immutableTagNS map[string]bool
	// Tag namespaces which are immutable on User and partially mutable on Topic:
	// user can only mutate tags he owns.
	maskedTagNS map[string]bool
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Sync()

	a := app{
		logger: l,
		utils:  utils.NewUtils(),
	}

	a.cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	a.statsChan = stats.NewStats(l)
	a.registerStatsVariables()

	server.Init(server.ServerArgs{
		Cfg:   a.cfg,
		Log:   a.logger,
		Stats: a.statsChan,
		Utils: a.utils,
	})

	if a.cfg.PProf.FileName != "" {
		if err := a.profile.StartProfile(a.cfg.PProf.FileName); err != nil {
			l.Sugar().Fatalf("failed to start profiling %w", err)
		}
	}

	a.initDBAdapter()
	defer func() {
		a.store.DBClose()
		a.logger.Info("Closed database connection(s)")
	}()

	a.initAuth()
	a.initValidators()
	a.initTags()

	if handChan := a.initMedia(); handChan != nil {
		defer func() {
			handChan <- true
			a.logger.Info("stopped files garbage collection")
		}()
	}

	if ch := a.initAccountGC(); ch != nil {
		defer func() {
			ch <- true
			a.logger.Info("stopped account garbage collector")
		}()
	}

	psh, err := push.NewPush(a.cfg.Push)
	if err != nil {
		a.logger.Fatal("failed to init push notifications", zapcore.Field{Interface: err})
	}

	defer func() {
		psh.Stop()
		a.logger.Info("stopped pushing notifications")
	}()

}
