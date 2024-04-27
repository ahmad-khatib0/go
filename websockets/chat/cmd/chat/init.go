package main

import (
	"expvar"
	"fmt"
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/handlers/files"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/users"
	"go.uber.org/zap/zapcore"
)

func (a *application) registerStatsVariables() {
	a.StatsChan.RegisterInt(constants.StatsVersion)
	decVer := a.Utils.Base10Version(a.Utils.ParseBuildstampVersion(a.Cfg.App.BuildStampCommand))
	if decVer <= 0 {
		decVer = a.Utils.Base10Version(a.Utils.ParseBuildstampVersion(a.Cfg.App.Version))
	}
	a.StatsChan.IntStatsSet(constants.StatsVersion, decVer)

	// Registering variables even if it's a standalone server. Otherwise
	// monitoring software will complain about missing vars.
	a.StatsChan.RegisterInt(constants.StatsClusterLeader)     // 1 if this node is cluster leader, 0 otherwise
	a.StatsChan.RegisterInt(constants.StatsClusterTotalNodes) // Total number of nodes configured
	a.StatsChan.RegisterInt(constants.StatsClusterLiveNodes)  // Number of nodes currently believed to be up.
}

func (a *application) initDBAdapter() {
	st, err := store.NewStore(store.StoreArgs{Logger: a.Logger})
	if err != nil {
		a.Logger.Fatal("failed to init store: %w", zapcore.Field{Interface: err})
	}

	a.Store = st
	a.Logger.Sugar().Infof("DB adapter: %s with version %d", a.Store.DBGetAdapterName(), a.Store.DBGetAdapterVersion())

	if f := a.Store.DBStats(); f != nil {
		expvar.Publish(constants.StatsDB, expvar.Func(f))
	}
}

func (a *application) initAuth() {
	err := a.Store.InitAuthLogicalNames(a.Cfg.Auth.LogicalNames)
	if err != nil {
		a.Logger.Sugar().Fatalf("failed to init auth %w: ", err)
	}

	// List of tag namespaces for user discovery which cannot be changed directly
	// by the client, e.g. 'email' or 'tel'.
	a.ImmutableTagNS = make(map[string]bool)

	authNames := a.Store.AuthGetAuthNames()
	for _, name := range authNames {

		if ah := a.Store.AuthGetLogicalAuthHandler(name); ah == nil {
			a.Logger.Sugar().Fatalf("unknown authenticator %s", ah)

		} else if jc := ah.GetAuthConfig(); jc != nil {

			if err := ah.Init(jc, name); err != nil {
				a.Logger.Sugar().Fatalf("failed to init auth scheme: %s, err: %w", name, err)
			}

			tags, err := ah.RestrictedTags()
			if err != nil {
				a.Logger.Sugar().Fatalf("failed get restricted tag namespaces (prefixes) for authenticator %s, %w", name, err)
			}

			for _, t := range tags {
				if strings.Contains(t, ":") {
					a.Logger.Sugar().Fatalf("tags restricted by auth handler should not contain character ':' %s", t)
				}
				a.ImmutableTagNS[name] = true
			}

		}
	}
}

func (a *application) initValidators() {
	type validator struct {
		Name      string
		AddToTags bool
		Required  []string
		Config    interface{}
	}

	validators := []validator{}
	if a.Cfg.Validator.Email != nil {
		e := a.Cfg.Validator.Email
		validators = append(validators, validator{Name: "email", AddToTags: e.AddToTags, Required: e.Required, Config: e})
	}

	for i, vc := range validators {
		name := validators[i].Name
		// Check if validator is restrictive. If so, add validator name to the list of restricted tags.
		// The namespace can be restricted even if the validator is disabled.
		if vc.AddToTags {
			if strings.Contains(name, ":") {
				a.Logger.Sugar().Fatalf("validator name should not contain  ':' character %s", name)
			}
			a.ImmutableTagNS[name] = true
		}

		if len(vc.Required) == 0 {
			// Skip disabled validator.  (i.e validating email is not required)
			continue
		}

		var rl []types.Level
		for _, r := range vc.Required {
			al := types.ParseAuthLevel(r)
			if al == types.LevelNone {
				a.Logger.Sugar().Fatalf("Invalid required AuthLevel '%s' in validator '%s'", r, name)
			}

			rl = append(rl, al)
			if a.AuthValidators == nil {
				a.AuthValidators = make(map[types.Level][]string)
			}

			a.AuthValidators[al] = append(a.AuthValidators[al], name)
		}

		if val := a.Store.GetValidator(name); val == nil {
			a.Logger.Fatal("Config provided for an unknown validator '" + name + "'")

		} else if err := val.Init(name, vc.Config); err != nil {
			a.Logger.Sugar().Fatalf("failed to init validator: %s, %w", name, err)
		}

		if a.Validators == nil {
			a.Validators = make(map[string]users.CredValidator)
		}

		a.Validators[name] = users.CredValidator{
			RequiredAuthLvl: rl,
			AddToTags:       vc.AddToTags,
		}
	}

	// Create credential validator config for clients.
	if len(a.AuthValidators) > 0 {
		a.ValidatorClientConfig = make(map[string][]string)
		for k, v := range a.AuthValidators {
			a.ValidatorClientConfig[k.String()] = v
		}
	}
}

func (a *application) initHandlers() {
	a.Handlers.Files = files.NewFilesHandler(a.Store.Adp(), a.Logger)
}

func (a *application) initTags() {
	// Partially restricted tag namespaces.
	a.MaskedTagNS = make(map[string]bool, len(a.Cfg.App.MaskedTagsNS))
	for _, t := range a.Cfg.App.MaskedTagsNS {
		if strings.Contains(t, ":") {
			a.Logger.Sugar().Fatalf("namespaces should not contain character -> ':'  for tag: %s", t)
		}
		a.MaskedTagNS[t] = true
	}

	var tags []string
	for t := range a.ImmutableTagNS {
		tags = append(tags, "'"+t+"'")
	}
	if len(tags) > 0 {
		a.Logger.Info("restricted tags: ", zapcore.Field{Interface: tags})
	}

	tags = nil
	for tag := range a.MaskedTagNS {
		tags = append(tags, "'"+tag+"'")
	}

	if len(tags) > 0 {
		a.Logger.Info("masked tags: ", zapcore.Field{Interface: tags})
	}
}

func (a *application) initMedia() chan<- bool {
	var ch chan<- bool

	if a.Cfg.Media != nil {

		if a.Cfg.Media.HandlerName == "" {
			a.Cfg.Media = nil
		} else {
			handlers := map[string]interface{}{}
			n := a.Cfg.Media.HandlerName

			if a.Cfg.Media.FS != nil {
				handlers["fs"] = a.Cfg.Media.FS
			}

			if err := a.Store.SetDefaultMediaHandler(n, handlers[n]); err != nil {
				a.Logger.Sugar().Fatalf("failed to init media handler %s, %w", n, err)
			}

			gp := a.Cfg.Media.GcPeriod
			gb := a.Cfg.Media.GcBlockSize

			if gp > 0 && gb > 0 {

				p, err := time.ParseDuration(fmt.Sprintf("%ds", gp))
				if err != nil {
					a.Logger.Sugar().Fatalf("failed to parse GcPeriod duration %w", err)
				}

				ch = a.Handlers.Files.InitLargeFilesGarbageCollection(p, gb)
			}
		}
	}

	return ch
}

func (a *application) initAccountGC() chan<- bool {
	var ch chan<- bool

	if a.Cfg.AccountGC != nil && a.Cfg.AccountGC.Enabled {
		gp := a.Cfg.AccountGC.GcPeriod
		gb := a.Cfg.AccountGC.GcBlockSize
		ga := a.Cfg.AccountGC.GcMinAccountAge

		if gp <= 0 || gb <= 0 || ga <= 0 {
			a.Logger.Fatal("invalid account gc config")
		}

		period := time.Second * time.Duration(gp)
		ch = a.users.InitUsersGarbageCollection(period, a.Cfg.AccountGC.GcBlockSize, a.Cfg.AccountGC.GcMinAccountAge)
	}

	return ch
}

func (a *application) initVideoCall() {
	if a.Cfg.Webrtc == nil || !a.Cfg.Webrtc.Enabled {
		return
	}

}
