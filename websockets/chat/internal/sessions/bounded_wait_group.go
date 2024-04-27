package session

import (
	"github.com/rs/zerolog/log"
)

func newBoundedWaitGroup(capacity int) *boundedWaitGroup {
	return &boundedWaitGroup{sem: make(chan struct{}, capacity)}
}

func (b *boundedWaitGroup) Add(delta int) {
	if delta <= 0 {
		return
	}

	for i := 0; i < delta; i++ {
		b.sem <- struct{}{}
	}

	b.wg.Add(delta)
}

func (b *boundedWaitGroup) Done() {
	select {
	case _, ok := <-b.sem:
		if !ok {
			log.Panic().Msg("boundedWaitGroup.sem closed.")
		}

	default:
		log.Panic().Msg("boundedWaitGroup.Done() called before Add().")
	}

	b.wg.Done()
}

func (b *boundedWaitGroup) Wait() {
	b.wg.Wait()
}
