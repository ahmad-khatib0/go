// Debug tooling. Dumps named profile in response to HTTP request at
//
//	http(s)://<host-name>/<configured-path>/<profile-name>
//
// See godoc for the list of possible profile names: https://golang.org/pkg/runtime/pprof/#Profile
package profile

import (
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
)

type Profile struct {
	Url string
}

func NewProfile(path string) *Profile {
	return &Profile{Url: path}
}

func (p *Profile) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	profileName := strings.TrimPrefix(r.URL.Path, p.Url)
	profile := pprof.Lookup(profileName)
	if profile == nil {
		p.err(w, http.StatusNotFound, "Unknown requested profile: "+profileName)
		return
	}

	profile.WriteTo(w, 2)
}

func (p *Profile) StartProfile(filepath string) error {
	cpuf, err := os.Create(filepath + ".cpu")
	if err != nil {
		return fmt.Errorf("failed to create cpu profiling file: %w", err)
	}
	defer cpuf.Close()

	memp, err := os.Create(filepath + ".mem")
	if err != nil {
		return fmt.Errorf("failed to create memory profiling file: %w", err)
	}
	defer memp.Close()

	pprof.StartCPUProfile(cpuf)
	defer pprof.StopCPUProfile()
	defer pprof.WriteHeapProfile(memp)
	return nil
}

func (p *Profile) err(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-GO-Pprof", "1")
	w.Header().Del("Content-Disposition")
	w.WriteHeader(code)

	fmt.Fprintln(w, msg)
}