package depot

import (
	"context"

	"github.com/ahmad-khatib0/go/event-driven-architecture/mallbots/internal/monolith"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono monolith.Monolith) error {

	return nil
}
