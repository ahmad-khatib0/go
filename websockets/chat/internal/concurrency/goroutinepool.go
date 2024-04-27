package concurrency

import "github.com/ahmad-khatib0/go/websockets/chat/internal/models"

// NewGoRoutinePool() allocates a new thread pool with `numWorkers` goroutines.
func NewGoRoutinePool(numWorkers int) models.GoRoutinePool {
	return goRoutinePool{
		work: make(chan models.Task),
		sem:  make(chan struct{}, numWorkers),
		stop: make(chan struct{}, numWorkers),
	}
}
