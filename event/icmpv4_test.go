package event_test

import (
	"context"
	"testing"

	"github.com/go-loadtest/evbundler/event"
)

func TestICMPv4(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ev := event.NewICMPv4Event(1, "0.0.0.0", "google.com")
	if err := ev.Fire(ctx); err != nil {
		t.Errorf("failed to send icmp from 0.0.0.0 to google.com: %v", err)
	}
}
