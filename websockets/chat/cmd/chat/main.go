package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/cluster"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/handlers/files"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/profile"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/push"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/users"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"go.uber.org/zap/zapcore"
)

type application struct {
	Logger                *logger.Logger
	Cfg                   *config.Config
	Store                 *store.Store
	StatsChan             *stats.Stats
	Utils                 *utils.Utils
	Cluster               models.Cluster
	Profile               *profile.Profile
	AuthValidators        map[types.Level][]string       // Validators required for each auth level
	Validators            map[string]users.CredValidator // Credential validators.
	ValidatorClientConfig map[string][]string            // Credential validator config to pass to clients.

	// Tag namespaces (prefixes) which are immutable to the client.
	ImmutableTagNS map[string]bool
	// Tag namespaces which are immutable on User and partially mutable on Topic:
	// user can only mutate tags he owns.
	MaskedTagNS map[string]bool
	users       users.Users

	Handlers struct {
		Files *files.FilesHandler
	}
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer l.Sync()

	a := application{
		Logger: l,
		Utils:  utils.NewUtils(),
	}

	a.Cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	a.StatsChan = stats.NewStats(l)
	a.registerStatsVariables()

	executable, _ := os.Executable()
	a.Logger.Info(fmt.Sprintf(
		"server: v%s%s%s pid: %d process(es): %d",
		a.Cfg.App.Version,
		executable,
		a.Cfg.App.BuildStampCommand,
		os.Getpid(),
		runtime.GOMAXPROCS(runtime.NumCPU()),
	))

	c, workerID, err := cluster.NewCluster(cluster.ClusterArgs{
		Cfg:    &a.Cfg.Cluster,
		Logger: l,
		Stats:  a.StatsChan,
	})

	if err != nil {
		l.Sugar().Fatalf("failed to init cluster %w", err)
	}
	a.Cluster = c

	if a.Cfg.PProf.FileName != "" {
		if err := a.Profile.StartProfile(a.Cfg.PProf.FileName); err != nil {
			l.Sugar().Fatalf("failed to start profiling %w", err)
		}
	}

	a.initDBAdapter()
	defer func() {
		a.Store.DBClose()
		a.Logger.Info("Closed database connection(s)")
	}()

	a.initAuth()
	a.initValidators()
	a.initHandlers()
	a.initTags()

	if handChan := a.initMedia(); handChan != nil {
		defer func() {
			handChan <- true
			a.Logger.Info("stopped files garbage collection")
		}()
	}

	if ch := a.initAccountGC(); ch != nil {
		defer func() {
			ch <- true
			a.Logger.Info("stopped account garbage collector")
		}()
	}

	psh, err := push.NewPush(a.Cfg.Push)
	if err != nil {
		a.Logger.Fatal("failed to init push notifications", zapcore.Field{Interface: err})
	}

	defer func() {
		psh.Stop()
		a.Logger.Info("stopped pushing notifications")
	}()

}
