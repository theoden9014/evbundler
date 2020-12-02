package event_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/theoden9014/evbundler/event"
)

func TestHTTPRequestEvent_Fire(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		method    string
		path      string
		body      io.Reader
		expectErr bool
	}{
		{"standard case", "GET", "/test", nil, false},
		{"invalid path", "GET", ":test", nil, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			sv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if got, want := req.Method, tt.method; got != want {
					t.Errorf("not expected method. got: %q , want: %q", got, want)
				}
				if got, want := req.URL.Path, tt.path; got != want {
					t.Errorf("not expected path. got: %q , want: %q", got, want)
				}
			}))
			defer sv.Close()

			addr := sv.Listener.Addr().String()
			u, err := url.Parse("http://" + addr)
			if err != nil {
				t.Fatal(err)
			}

			ev := &event.HTTPEvent{
				URL:    u,
				Method: tt.method,
				Path:   tt.path,
			}

			err = ev.Fire(ctx)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expect raise error, but not raise error in HTTPEvent: %+v", ev)
				}
			} else {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
