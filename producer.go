package evbundler

import (
	"context"
	"time"
)

type tickerProducer struct {
	ticker        *time.Ticker
	sendEventFunc func(time.Time, chan<- Event)
}

func TickerProducer(ctx context.Context, d time.Duration, f func(time.Time, chan<- Event)) chan Event {
	p := tickerProducer{
		ticker:        time.NewTicker(d),
		sendEventFunc: f,
	}
	return p.Produce(ctx)
}

func (p *tickerProducer) Produce(ctx context.Context) chan Event {
	evCh := make(chan Event)

	go func() {
		defer p.ticker.Stop()

		for {
			select {
			case t := <-p.ticker.C:
				p.sendEventFunc(t, evCh)
			case <-ctx.Done():
				close(evCh)
				return
			}
		}
	}()

	return evCh
}
