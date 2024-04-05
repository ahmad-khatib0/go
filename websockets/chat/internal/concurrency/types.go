package concurrency

// Task represents a work task to be run on the specified thread pool.
type Task func()

// GoRoutinePool is a pull of Go routines with associated locking mechanism.
type GoRoutinePool struct {
	// Work queue.
	work chan Task
	// Counter to control the number of already allocated/running goroutines.
	sem chan struct{}
	// Exit knob.
	stop chan struct{}
}
