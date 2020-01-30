package event

import (
	"context"
	"time"
)

type Trigger interface {
	Channel(context.Context) chan Event
}

type tickerTrigger struct {
	ticker        *time.Ticker
	sendEventFunc func(time.Time, chan Event)
}

func TickerTrigger(ctx context.Context, d time.Duration, f func(time.Time, chan Event)) chan Event {
	g := tickerTrigger{
		ticker:        time.NewTicker(d),
		sendEventFunc: f,
	}
	return g.Channel(ctx)
}

func (g *tickerTrigger) Channel(ctx context.Context) chan Event {
	eventChan := make(chan Event)

	go func() {
		defer g.ticker.Stop()

		for {
			select {
			case t := <-g.ticker.C:
				g.sendEventFunc(t, eventChan)
			case <-ctx.Done():
				close(eventChan)
				return
			}
		}
	}()

	return eventChan
}
