package registry

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aplulu/etcd-shim/internal/driver"
)

type NewDriverFn func(ctx context.Context, log *slog.Logger) (driver.Driver, error)

var (
	ErrDriverNotFound = errors.New("driver not found")

	driverRegistry = map[string]NewDriverFn{}
)

func Register(name string, fn NewDriverFn) {
	driverRegistry[name] = fn
}

func NewDriver(name string, ctx context.Context, log *slog.Logger) (driver.Driver, error) {
	fn, ok := driverRegistry[name]
	if !ok {
		return nil, ErrDriverNotFound
	}
	return fn(ctx, log)
}
