package concurrency

import "github.com/ahmad-khatib0/go/websockets/chat/internal/models"

type goRoutinePool struct {
	// Work queue.
	work chan models.Task
	// Counter to control the number of already allocated/running goroutines.
	sem chan struct{}
	// Exit knob.
	stop chan struct{}
}
