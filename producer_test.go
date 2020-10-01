package evbundler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-loadtest/evbundler"
)

type errEvent struct {
	err error
}

func (f errEvent) Name() string                   { return "err" }
func (f errEvent) Fire(ctx context.Context) error { return f.err }

func TestTickerProducer(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
	defer cancel()

	var fooErr = errors.New("foo")
	evch := evbundler.TickerProducer(ctx, 1*time.Millisecond, func(t time.Time, evch chan<- evbundler.Event) {
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
