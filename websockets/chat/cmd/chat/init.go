package main

import (
	"expvar"
	"fmt"
	"strings"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/server"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"go.uber.org/zap/zapcore"
)

func (a *app) registerStatsVariables() {
	a.statsChan.RegisterInt(constants.StatsVersion)
	decVer := a.utils.Base10Version(a.utils.ParseBuildstampVersion(a.cfg.App.BuildStampCommand))
	if decVer <= 0 {
		decVer = a.utils.Base10Version(a.utils.ParseBuildstampVersion(a.cfg.App.Version))
	}
	a.statsChan.IntStatsSet(constants.StatsVersion, decVer)

	// Registering variables even if it's a standalone server. Otherwise
	// monitoring software will complain about missing vars.
	a.statsChan.RegisterInt(constants.StatsClusterLeader)     // 1 if this node is cluster leader, 0 otherwise
	a.statsChan.RegisterInt(constants.StatsClusterTotalNodes) // Total number of nodes configured
	a.statsChan.RegisterInt(constants.StatsClusterLiveNodes)  // Number of nodes currently believed to be up.
}

func (a *app) initDBAdapter() {
	st, err := store.NewStore(store.StoreArgs{Logger: a.logger})
	if err != nil {
		a.logger.Fatal("failed to init store: %w", zapcore.Field{Interface: err})
	}

	a.store = st
	a.logger.Sugar().Infof("DB adapter: %s with version %d", a.store.DBGetAdapterName(), a.store.DBGetAdapterVersion())

	if f := a.store.DBStats(); f != nil {
		expvar.Publish(constants.StatsDB, expvar.Func(f))
	}
}

func (a *app) initAuth() {
	err := a.store.InitAuthLogicalNames(a.cfg.Auth.LogicalNames)
	if err != nil {
		a.logger.Sugar().Fatalf("failed to init auth %w: ", err)
	}

	// List of tag namespaces for user discovery which cannot be changed directly
	// by the client, e.g. 'email' or 'tel'.
	a.immutableTagNS = make(map[string]bool)

	authNames := a.store.AuthGetAuthNames()
	for _, name := range authNames {

		if ah := a.store.AuthGetLogicalAuthHandler(name); ah == nil {
			a.logger.Sugar().Fatalf("unknown authenticator %s", ah)

		} else if jc := ah.GetAuthConfig(); jc != nil {

			if err := ah.Init(jc, name); err != nil {
				a.logger.Sugar().Fatalf("failed to init auth scheme: %s, err: %w", name, err)
			}

			tags, err := ah.RestrictedTags()
			if err != nil {
				a.logger.Sugar().Fatalf("failed get restricted tag namespaces (prefixes) for authenticator %s, %w", name, err)
			}

			for _, t := range tags {
				if strings.Contains(t, ":") {
					a.logger.Sugar().Fatalf("tags restricted by auth handler should not contain character ':' %s", t)
				}
				a.immutableTagNS[name] = true
			}

		}
	}
}

func (a *app) initValidators() {
	type validator struct {
		Name      string
		AddToTags bool
		Required  []string
		Config    interface{}
	}

	validators := []validator{}
	if a.cfg.Validator.Email != nil {
		e := a.cfg.Validator.Email
		validators = append(validators, validator{Name: "email", AddToTags: e.AddToTags, Required: e.Required, Config: e})
	}

	for i, vc := range validators {
		name := validators[i].Name
		// Check if validator is restrictive. If so, add validator name to the list of restricted tags.
		// The namespace can be restricted even if the validator is disabled.
		if vc.AddToTags {
			if strings.Contains(name, ":") {
				a.logger.Sugar().Fatalf("validator name should not contain  ':' character %s", name)
			}
			a.immutableTagNS[name] = true
		}

		if len(vc.Required) == 0 {
			// Skip disabled validator.  (i.e validating email is not required)
			continue
		}

		var rl []types.Level
		for _, r := range vc.Required {
			al := types.ParseAuthLevel(r)
			if al == types.LevelNone {
				a.logger.Sugar().Fatalf("Invalid required AuthLevel '%s' in validator '%s'", r, name)
			}

			rl = append(rl, al)
			if a.authValidators == nil {
				a.authValidators = make(map[types.Level][]string)
			}

			a.authValidators[al] = append(a.authValidators[al], name)
		}

		if val := a.store.GetValidator(name); val == nil {
			a.logger.Fatal("Config provided for an unknown validator '" + name + "'")

		} else if err := val.Init(name, vc.Config); err != nil {
			a.logger.Sugar().Fatalf("failed to init validator: %s, %w", name, err)
		}

		if a.validators == nil {
			a.validators = make(map[string]server.CredValidator)
		}

		a.validators[name] = server.CredValidator{
			RequiredAuthLvl: rl,
			AddToTags:       vc.AddToTags,
		}
	}

	// Create credential validator config for clients.
	if len(a.authValidators) > 0 {
		a.validatorCliCfg = make(map[string][]string)
		for k, v := range a.authValidators {
			a.validatorCliCfg[k.String()] = v
		}
	}
}

func (a *app) initTags() {
	// Partially restricted tag namespaces.
	a.maskedTagNS = make(map[string]bool, len(a.cfg.App.MaskedTagsNS))
	for _, t := range a.cfg.App.MaskedTagsNS {
		if strings.Contains(t, ":") {
			a.logger.Sugar().Fatalf("namespaces should not contain character -> ':'  for tag: %s", t)
		}
		a.maskedTagNS[t] = true
	}

	var tags []string
	for t := range a.immutableTagNS {
		tags = append(tags, "'"+t+"'")
	}
	if len(tags) > 0 {
		a.logger.Info("restricted tags: ", zapcore.Field{Interface: tags})
	}

	tags = nil
	for tag := range a.maskedTagNS {
		tags = append(tags, "'"+tag+"'")
	}

	if len(tags) > 0 {
		a.logger.Info("masked tags: ", zapcore.Field{Interface: tags})
	}
}

func (a *app) initMedia() chan<- bool {
	var ch chan<- bool

	if a.cfg.Media != nil {

		if a.cfg.Media.HandlerName == "" {
			a.cfg.Media = nil
		} else {
			handlers := map[string]interface{}{}
			n := a.cfg.Media.HandlerName

			if a.cfg.Media.FS != nil {
				handlers["fs"] = a.cfg.Media.FS
			}

			if err := a.store.SetDefaultMediaHandler(n, handlers[n]); err != nil {
				a.logger.Sugar().Fatalf("failed to init media handler %s, %w", n, err)
			}

			gp := a.cfg.Media.GcPeriod
			gb := a.cfg.Media.GcBlockSize

			if gp > 0 && gb > 0 {

				p, err := time.ParseDuration(fmt.Sprintf("%ds", gp))
				if err != nil {
					a.logger.Sugar().Fatalf("failed to parse GcPeriod duration %w", err)
				}

				ch = server.InitUsersGarbageCollection(p, gb, a.cfg.AccountGC.GcMinAccountAge)
			}
		}
	}

	return ch
}

func (a *app) initAccountGC() chan<- bool {
	var ch chan<- bool

	if a.cfg.AccountGC != nil && a.cfg.AccountGC.Enabled {
		gp := a.cfg.AccountGC.GcPeriod
		gb := a.cfg.AccountGC.GcBlockSize
		ga := a.cfg.AccountGC.GcMinAccountAge

		if gp <= 0 || gb <= 0 || ga <= 0 {
			a.logger.Fatal("invalid account gc config")
		}

		period := time.Second * time.Duration(gp)
		ch = server.InitUsersGarbageCollection(period, a.cfg.AccountGC.GcBlockSize, a.cfg.AccountGC.GcMinAccountAge)
	}

	return ch
}

func (a *app) initVideoCall() {
	if a.cfg.Webrtc == nil || !a.cfg.Webrtc.Enabled {
		return
	}

}
