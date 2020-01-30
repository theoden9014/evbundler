package event

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var HTTPClient httpClient = http.DefaultClient

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPEvent struct {
	URL    *url.URL
	Method string
	Path   string
	Body   io.Reader
}

func (he HTTPEvent) Name() string {
	return "HTTP"
}

func (he HTTPEvent) Fire(ctx context.Context) error {
	u, err := he.URL.Parse(he.Path)
	if err != nil {
		return fmt.Errorf("invalid URL Scheme: %v, or ivalid Path: %v : %w", he.URL.String(), he.Path, err)
	}
	req, err := http.NewRequestWithContext(ctx, he.Method, u.String(), he.Body)
	if err != nil {
		return fmt.Errorf("can not create new http.Request: %w", err)
	}

	// should *http.Client to be in HTTPEvent?
	resp, err := HTTPClient.Do(req)

	if err != nil {
		return fmt.Errorf("can not send a HTTP request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New("unsuccessful response code.")
	}

	return nil
}
