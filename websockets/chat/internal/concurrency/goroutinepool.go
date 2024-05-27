package concurrency

import "github.com/ahmad-khatib0/go/websockets/chat/internal/models"

// NewGoRoutinePool() allocates a new thread pool with `numWorkers` goroutines.
func NewGoRoutinePool(numWorkers int) models.GoRoutinePool {
	return &goRoutinePool{
		work: make(chan models.Task),
		sem:  make(chan struct{}, numWorkers),
		stop: make(chan struct{}, numWorkers),
	}
}

// Schedule enqueus a closure to run on the GoRoutinePool's goroutines.
func (gr *goRoutinePool) Schedule(t models.Task) {
	select {
	case gr.work <- t:
	case gr.sem <- struct{}{}:
		go gr.worker(t)
	}
}

// Stop sends a stop signal to all running goroutines.
func (p *goRoutinePool) Stop() {
	numWorkers := cap(p.sem)
	for i := 0; i < numWorkers; i++ {
		p.stop <- struct{}{}
	}
}

// Thread pool worker goroutine.
func (gr *goRoutinePool) worker(t models.Task) {
	defer func() {
		<-gr.sem
	}()

	for {
		t()
		select {
		case t = <-gr.work:
		case <-gr.sem:
			return
		}
	}
}
