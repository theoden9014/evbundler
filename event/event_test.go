package event_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theoden9014/evbundler/event"
)

type errEvent struct {
	err error
}

func (f errEvent) Name() string                   { return "err" }
func (f errEvent) Fire(ctx context.Context) error { return f.err }

func TestTickerEventGenerator(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()

	var fooErr = errors.New("foo")
	evch := event.TickerTrigeger(ctx, 1*time.Millisecond, func(t time.Time, evch chan event.Event) {
		evch <- &errEvent{err: fooErr}
	})

	for i := 0; i < 10; i++ {
		select {
		case ev := <-evch:
			err := ev.Fire(ctx)
			if !errors.As(fooErr, &err) {
				t.Error("no match error type")
			}
		case <-ctx.Done():
			t.Errorf("not enough to generate event in time: (count: %d/10)", i)
		}
	}
}
