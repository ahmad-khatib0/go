package main

import (
	"expvar"
	"net/http"
	"path"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/profile"
	"github.com/go-chi/chi/v5"
)

func (a *app) runMux() http.Handler {
	mux := chi.NewRouter()

	if a.cfg.Paths.Expvar != "" {
		mux.Get(a.cfg.Paths.Expvar, expvar.Handler().ServeHTTP)
		a.logger.Sugar().Infof("stats: variables exposed at: '%s' ", a.cfg.Paths.Expvar)
	}

	if a.cfg.Paths.PProf != "" {
		a.profile = profile.NewProfile(a.cfg.Paths.PProf)
		mux.Get(path.Clean("/"+a.cfg.Paths.PProf+"/"), a.profile.ProfileHandler)
		a.logger.Sugar().Infof("pprof: profiling info exposed at '%s'", a.cfg.Paths.PProf)
	}

	return mux
}
