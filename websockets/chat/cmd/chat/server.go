package main

import (
	"expvar"
	"net/http"
	"path"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/profile"
	"github.com/go-chi/chi/v5"
)

func (a *application) runMux() http.Handler {
	mux := chi.NewRouter()

	if a.Cfg.Paths.Expvar != "" {
		mux.Get(a.Cfg.Paths.Expvar, expvar.Handler().ServeHTTP)
		a.Logger.Sugar().Infof("stats: variables exposed at: '%s' ", a.Cfg.Paths.Expvar)
	}

	if a.Cfg.Paths.PProf != "" {
		a.Profile = profile.NewProfile(a.Cfg.Paths.PProf)
		mux.Get(path.Clean("/"+a.Cfg.Paths.PProf+"/"), a.Profile.ProfileHandler)
		a.Logger.Sugar().Infof("pprof: profiling info exposed at '%s'", a.Cfg.Paths.PProf)
	}

	return mux
}
