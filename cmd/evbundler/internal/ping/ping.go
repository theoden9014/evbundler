package ping

import (
	"context"
	"flag"
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
	Run:  run,
}

var target string
var interval string

func run(args []string) error {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.StringVar(&target, "target", "", `send to host or ip`)
	flags.StringVar(&interval, "interval", "1s", `interval of send each packets e.g. 1s, 300ms, 1m`)

	if err := flags.Parse(args); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d, err := time.ParseDuration(interval)
	if err != nil {
		return err
	}
	evCh := evbundler.TickerProducer(ctx, d, func(t time.Time, evCh chan evbundler.Event) {
		evCh <- event.NewICMPv4Event(1, "0.0.0.0", target)
	})

	wp := evbundler.NewWorkerPool(10, nil)
	disp := dispatcher.NewGoChannel(wp)
	_ = disp.Dispatch(ctx, evCh)

	return nil
}
