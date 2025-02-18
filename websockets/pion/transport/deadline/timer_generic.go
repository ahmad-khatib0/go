//go:build !js
// +build !js

package deadline

import "time"

func afterFunc(d time.Duration, f func()) timer {
	return time.AfterFunc(d, f)
}
