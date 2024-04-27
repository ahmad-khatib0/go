package files

import (
	"math/rand"
	"time"
)

// InitLargeFilesGarbageCollection() runs every 'period' and deletes up to 'blockSize'
//
// unused files. Returns a writable channel which can be used to stop the process.
func (s *FilesHandler) InitLargeFilesGarbageCollection(period time.Duration, blockSize int) chan<- bool {
	// Unbuffered stop channel. Whomever stops the gc must wait for the process to finish.
	stop := make(chan bool)

	go func() {
		// Add some randomness to the tick period to desynchronize runs on cluster nodes:
		// 0.75 * period + rand(0, 0.5) * period.
		period = (period >> 1) + (period >> 2) + time.Duration(rand.Intn(int(period>>1)))
		gcTicker := time.Tick(period)

		for {
			select {
			case <-gcTicker:
				deletedRecords, err := s.db.Files().DeleteUnused(time.Now().Add(-time.Hour), blockSize)
				if err != nil {
					s.logger.Sugar().Warnf("error deleting unused larage file from db <InitLargeFilesGarabageCollection> %w", err)
				} else if len(deletedRecords) > 0 {
					s.logger.Sugar().Warnf("deleting media files: %v", deletedRecords)
					err := s.media.Delete(deletedRecords)
					if err != nil {
						s.logger.Sugar().Warnf("error deleting unused larage file from storage <InitLargeFilesGarabageCollection> %w", err)
					}
				}

			case <-stop:
				return
			}
		}
	}()

	return stop
}
