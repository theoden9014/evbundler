package evbundler

import "time"

type Result struct {
	WorkerID  string        `json:"worker_id"`
	Weight    int           `json:"weight"`
	Error     error         `json:"error"`
	Latency   time.Duration `json:"latency"`
	Timestamp time.Time     `json:"timestamp"`
}

func (r *Result) End() time.Time { return r.Timestamp.Add(r.Latency) }

type Results []Result

func (rs *Results) Add(r *Result) { *rs = append(*rs, *r) }

func (rs Results) Len() int           { return len(rs) }
func (rs Results) Less(i, j int) bool { return rs[i].Timestamp.Before(rs[j].Timestamp) }
func (rs Results) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
