package packetio

import "errors"

var (
	// ErrFull is returned when the buffer has hit the configured limits.
	ErrFull = errors.New("packetio.Buffer is full, discarding write")

	// ErrTimeout is returned when a deadline has expired.
	ErrTimeout = errors.New("i/o timeout")
)

// netError implements net.Error.
type netError struct {
	error
	timeout, temporary bool
}

func (e *netError) Timeout() bool {
	return e.timeout
}

func (e *netError) Temporary() bool {
	return e.temporary
}
