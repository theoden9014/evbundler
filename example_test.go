package evbundler_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"time"

	"github.com/theoden9014/evbundler"
	"github.com/theoden9014/evbundler/dispatcher"
	"github.com/theoden9014/evbundler/event"
)

func ExampleTickerProducer() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	countupCh := make(chan int, 10)
	go func() {
		var counter int
		for {
			select {
			case <-countupCh:
				counter++
			default:
				if counter >= 100 {
					fmt.Print("receive 100 requests\n")
					cancel()
					return
				}
			}
		}
	}()

	sv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		countupCh <- 1
	}))
	addr := sv.Listener.Addr().String()
	defer sv.Close()

	u, err := url.Parse("http://" + addr)
	if err != nil {
		log.Fatal(err)
	}

	charset := "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// generate http events with random path.
	evCh := evbundler.TickerProducer(ctx, 1*time.Millisecond, func(t time.Time, evCh chan evbundler.Event) {
		var sb strings.Builder
		r := rand.New(rand.NewSource(t.UnixNano()))
		for i := 0; i < 10; i++ {
			sb.WriteByte(charset[r.Intn(len(charset))])
		}

		evCh <- event.HTTPEvent{
			URL:    u,
			Method: "GET",
			Path:   sb.String(),
			Body:   nil,
		}
	})

	wp := evbundler.NewWorkerPool(10, nil)
	disp := dispatcher.NewGoChannel(wp)
	_ = disp.Dispatch(ctx, evCh)

	// Output: receive 100 requests
}
