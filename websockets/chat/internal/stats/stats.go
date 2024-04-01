package stats

import (
	"expvar"
	"fmt"
	"runtime"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/histogram"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/models"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type Stats struct {
	ch     chan *models.StatsChanVariable
	Logger *logger.Logger
}

func NewStats(l *logger.Logger) *Stats {
	return &Stats{Logger: l}
}

func (s *Stats) RunStat() {
	s.ch = make(chan *models.StatsChanVariable, 1024)
	t := time.Now()

	expvar.Publish(constants.StatsUptime, expvar.Func(func() any {
		return time.Since(t).Seconds()
	}))

	expvar.Publish(constants.StatsNumGoroutines, expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	go s.statsUpdater()

}

// statsUpdater() using go routing to publish statistics updates
func (s *Stats) statsUpdater() {
	for up := range s.ch {
		if up == nil {
			s.ch = nil
			break
		}

		if v := expvar.Get(up.Varname); v != nil {
			switch val := v.(type) {
			case *expvar.Int:
				count := up.Value.(int64)
				if up.Inc {
					val.Add(count)
				} else {
					val.Set(count)
				}
			case *histogram.Histogram:
				value := up.Value.(float64)
				val.AddSample(value)
			default:
				s.Logger.Fatal(fmt.Sprintf("unsupported expvar type %T", val))
			}
		} else {
			s.Logger.Fatal(fmt.Sprintf("you are updating an unknown statistics variable!"))
		}
	}

	s.Logger.Info("stats: shutdown")
}

// RegisterInt() Registers integer variable. without checking for initialization.
func (s *Stats) RegisterInt(vn string) {
	expvar.Publish(vn, new(expvar.Int))
}

// Register histogram variable. `bounds` specifies histogram buckets/bins
// (see comment next to the `histogram` struct definition).
func (s *Stats) RegisterHistogram(vn string, bounds []float64) {
	numBuckets := len(bounds) + 1
	expvar.Publish(vn, &histogram.Histogram{
		CountPerBucket: make([]int64, numBuckets),
		Bounds:         bounds,
	})
}

// Async publish int variable
func (s *Stats) IntStatsSet(vn string, val int64) {
	if s.ch != nil {
		select {
		case s.ch <- &models.StatsChanVariable{Varname: vn, Value: val, Inc: false}:
		}
	}
}

// IntStatsInc() async publish an increment/decrement to int variable
func (s *Stats) IntStatsInc(vn string, val int) {
	if s.ch != nil {
		select {
		case s.ch <- &models.StatsChanVariable{Varname: vn, Value: int64(val), Inc: true}:
		}
	}
}

// Async publish a value (add a sample) to a histogram variable.
func (s *Stats) HistogramAddSample(vn string, val float64) {
	if s.ch != nil {
		select {
		case s.ch <- &models.StatsChanVariable{Varname: vn, Value: val}:
		default:
		}
	}
}

// Shutdown() stop publishing stats
func (s *Stats) Shutdown() {
	if s.ch != nil {
		s.ch <- nil
	}
}
