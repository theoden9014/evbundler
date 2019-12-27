package event

import (
	"io"
	"net/url"
)

type HTTPEvent struct {
	url    *url.URL
	method string
	path   string
	body   io.Reader
}
