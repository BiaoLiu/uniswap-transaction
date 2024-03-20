package alerting

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"

	"uniswap-transaction/pkg/alert"
)

// Server is an server alerting middleware.
func Server(alertManager *alert.Manager) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			alertManager.SendServiceMessage(ctx, req, err)
			return reply, err
		}
	}
}
