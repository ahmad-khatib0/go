package users

import (
	"math/rand"
	"time"

	"go.uber.org/zap/zapcore"
)

// InitUsersGarbageCollection() runs every 'period' and deletes up to 'blockSize'
//
// stale unvalidated user accounts which have been last updated at least 'minAccountAgeHours' hours.
//
// Returns channel which can be used to stop the process.
func (u *Users) InitUsersGarbageCollection(period time.Duration, blockSize, minAccountAgeHours int) chan<- bool {
	// Unbuffered stop channel. Whomever stops the gc must wait for the process to finish.
	stop := make(chan bool)

	go func() {
		// Add some randomness to the tick period to desynchronize runs on cluster nodes:
		// 0.75 * period + rand(0, 0.5) * period.
		period = period - (period >> 2) + time.Duration(rand.Intn(int(period>>1)))
		gt := time.Tick(period)

		u.logger.Sugar().Infof(
			"stale account GC started with period %s, block size %d, min account age %d hours",
			period.Round(time.Second),
			blockSize,
			minAccountAgeHours,
		)

		staleAge := time.Hour * time.Duration(minAccountAgeHours)

		for {
			select {
			case <-gt:
				if uids, err := u.db.Users().GetUnvalidated(time.Now().Add(-staleAge), blockSize); err == nil {

					if len(uids) > 0 {
						u.logger.Info("stale account GC will delete these uids: ", zapcore.Field{Interface: uids})
						for _, id := range uids {
							if err := u.db.Users().Delete(id, true); err != nil {
								u.logger.Sugar().Warnf("stale account GC failed to delete %s: %+v", id.UserId(), err)
							}
						}
					}

				} else {
					u.logger.Sugar().Infof("stale account GC error %w", err)
				}

			case <-stop:
				return
			}
		}

	}()

	return stop
}
