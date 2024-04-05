package concurrency

// NewGoRoutinePool() allocates a new thread pool with `numWorkers` goroutines.
func NewGoRoutinePool(numWorkers int) *GoRoutinePool {
	return &GoRoutinePool{
		work: make(chan Task),
		sem:  make(chan struct{}, numWorkers),
		stop: make(chan struct{}, numWorkers),
	}
}
