package event

import (
	"context"
)

type Event interface {
	Name() string
	Fire(ctx context.Context) error
}
