package ping

import (
	"context"
	"time"

	"github.com/go-loadtest/evbundler"
	"github.com/go-loadtest/evbundler/cmd/evbundler/internal/base"
	"github.com/go-loadtest/evbundler/dispatcher"
	"github.com/go-loadtest/evbundler/event"
)

const name = "ping"

var Cmd = &base.Command{
	Name: name,
	Doc:  "icmpv4 event generator",
}

func init() {
	Cmd.Run = run
}

var (
	target   = Cmd.Flags.String("target", "", `send to host or ip`)
	interval = Cmd.Flags.String("interval", "1s", `interval of send each packets e.g. 1s, 300ms, 1m`)
)

func run(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d, err := time.ParseDuration(*interval)
	if err != nil {
		return err
	}
	evCh := evbundler.TickerProducer(ctx, d, func(t time.Time, evCh chan evbundler.Event) {
		evCh <- event.NewICMPv4Event(1, "0.0.0.0", *target)
	})

	wp := evbundler.NewWorkerPool(10, nil)
	disp := dispatcher.NewGoChannel(wp)
	_ = disp.Dispatch(ctx, evCh)

	return nil
}
