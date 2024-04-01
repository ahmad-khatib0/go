package histogram

import (
	"encoding/json"
	"sort"
)

//	Histogram is a simple implementation of histogram expvar.Var.
//
// `Bounds` specifies the histogram buckets as follows (length = len(bounds)):
//
//	(-inf, Bounds[i]) for i = 0
//	[Bounds[i-1], Bounds[i]) for 0 < i < length
//	[Bounds[i-1], +inf) for i = length
type Histogram struct {
	Count          int64     `json:"count"`
	Sum            float64   `json:"sum"`
	CountPerBucket []int64   `json:"count_per_bucket"`
	Bounds         []float64 `json:"bounds"`
}

func (h *Histogram) AddSample(v float64) {
	h.Count++
	h.Sum += v
	idx := sort.SearchFloat64s(h.Bounds, v)
	h.CountPerBucket[idx]++
}

func (h *Histogram) String() string {
	if r, err := json.Marshal(h); err == nil {
		return string(r)
	}
	return ""
}
