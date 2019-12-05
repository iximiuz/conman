package timeutil

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

func WithTimeout(d time.Duration, fn func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	ch := make(chan error, 1)

	func() {
		ch <- fn()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), "Timed out")
	}
}
