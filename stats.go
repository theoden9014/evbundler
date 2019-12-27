package evbundler

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// The following Event measures are supported for use in custom views.
var (
	EventCount = stats.Int64("evbundler/event_count", "Number of processed events", stats.UnitDimensionless)
)

var (
	KeyEventErr, _ = tag.NewKey("event_err")
)

// The following Workers measures are supported for use in custom views.
var (
	WorkerCount = stats.Int64("evbundler/worker_count", "Number of total workers", stats.UnitDimensionless)
)

var (
	// WorkerStatus is the status of worker, capitalized (i.e. DEAD, ACTIVE, WAIT_EVENT, PROCESS_EVENT)
	// See also `type WorkerState`
	KeyWorkerState, _ = tag.NewKey("worker.state")
)
