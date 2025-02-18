package deadline

import "time"

type timer interface {
	Stop() bool
	Reset(time.Duration) bool
}
