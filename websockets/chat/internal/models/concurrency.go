package models

// Task represents a work task to be run on the specified thread pool.
type Task func()

// GoRoutinePool is a pull of Go routines with associated locking mechanism.
type GoRoutinePool interface{}
