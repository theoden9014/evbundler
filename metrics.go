package evbundler

import (
	"encoding/json"
	"time"

	"github.com/influxdata/tdigest"
)

func NewMetrics() *metrics {
	return &metrics{}
}

type metrics struct {
	Latencies LatencyMetrics `json:"latencies"`
	Earliest  time.Time      `json:"-"`
	Latest    time.Time      `json:"-"`
	End       time.Time      `json:"-"`

	Events  uint64            `json:"events"`
	Success uint64            `json:"success"`
	Failed  uint64            `json:"failed"`
	Errors  map[string]uint64 `json:"errors"`

	Rate       float64       `json:"rate"`
	Throughput float64       `json:"throughput"`
	Duration   time.Duration `json:"duration"`
}

func (m *metrics) Add(r *Result) {
	m.init()

	m.Latencies.Add(r.Latency)
	m.Events += uint64(r.Weight)

	if m.Earliest.IsZero() || m.Earliest.After(r.Timestamp) {
		m.Earliest = r.Timestamp
	}

	if r.Timestamp.After(m.Latest) {
		m.Latest = r.Timestamp
	}

	if end := r.End(); end.After(m.End) {
		m.End = end
	}

	if r.Error != nil {
		m.Failed += uint64(r.Weight)
		err := r.Error.Error()
		m.Errors[err]++
	} else {
		m.Success += uint64(r.Weight)
	}
}

func (m *metrics) Export() metrics {
	m2 := *m
	m2.init()

	m2.Rate = float64(m2.Events)
	m2.Throughput = float64(m2.Success)
	m2.Duration = m2.Latest.Sub(m2.Earliest)

	if mins := m2.Duration.Minutes(); mins > 0 {
		m2.Rate /= mins
		m2.Throughput /= mins
	}

	m2.Latencies.Mean = time.Duration(float64(m2.Latencies.Total) / float64(m2.Events))
	m2.Latencies.P50 = m.Latencies.Quantile(0.50)
	m2.Latencies.P90 = m.Latencies.Quantile(0.90)
	m2.Latencies.P95 = m.Latencies.Quantile(0.95)
	m2.Latencies.P99 = m.Latencies.Quantile(0.99)

	return m2
}

func (m *metrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *metrics) init() {
	if m.Errors == nil {
		m.Errors = make(map[string]uint64)
	}
}

type LatencyMetrics struct {
	Total time.Duration
	Mean  time.Duration
	P50   time.Duration
	P90   time.Duration
	P95   time.Duration
	P99   time.Duration
	Max   time.Duration

	estimator estimator
}

func (l *LatencyMetrics) Add(latency time.Duration) {
	l.init()
	if l.Total += latency; latency > l.Max {
		l.Max = latency
	}
	l.estimator.Add(float64(latency))
}

func (l LatencyMetrics) Quantile(nth float64) time.Duration {
	l.init()
	return time.Duration(l.estimator.Get(nth))
}

func (l *LatencyMetrics) init() {
	if l.estimator == nil {
		l.estimator = newTdigestEstimator(100)
	}
}

func (l *LatencyMetrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Total string `json:"-"`
		Mean  string `json:"mean"`
		P50   string `json:"50th"`
		P90   string `json:"90th"`
		P95   string `json:"95th"`
		P99   string `json:"99th"`
		Max   string `json:"max"`
	}{
		Total: l.Total.String(),
		Mean:  l.Mean.String(),
		P50:   l.P50.String(),
		P90:   l.P90.String(),
		P95:   l.P95.String(),
		P99:   l.P99.String(),
		Max:   l.Max.String(),
	})
}

type estimator interface {
	Add(sample float64)
	Get(quantile float64) float64
}

type tdigestEstimator struct{ *tdigest.TDigest }

func newTdigestEstimator(compression float64) *tdigestEstimator {
	return &tdigestEstimator{TDigest: tdigest.NewWithCompression(compression)}
}

func (e *tdigestEstimator) Add(s float64) { e.TDigest.Add(s, 1) }
func (e *tdigestEstimator) Get(q float64) float64 {
	return e.TDigest.Quantile(q)
}
