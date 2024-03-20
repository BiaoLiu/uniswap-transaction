package pprof

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

const (
	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrForbidden = errors.Forbidden(reason, "Forbidden")
)

// Option is pprof option.
type Option func(*options)

// Parser is a jwt parser
type options struct {
	switchHandler func() bool
}

// WithSwitchHandler with switch handler option.
func WithSwitchHandler(f func() bool) Option {
	return func(o *options) {
		o.switchHandler = f
	}
}

// Server is a pprof middleware.
func Server(opts ...Option) middleware.Middleware {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if o.switchHandler != nil {
				pprofSwitch := o.switchHandler()
				if !pprofSwitch {
					return nil, ErrForbidden
				}
			}
			return handler(ctx, req)
		}
	}
}
