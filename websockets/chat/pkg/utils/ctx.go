package utils

import (
	"context"
	"time"
)

func (u *Utils) GetContext(dur time.Duration) (context.Context, context.CancelFunc) {
	if dur > 0 {
		return context.WithTimeout(context.Background(), dur)
	}
	return context.Background(), nil
}
