package alert

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/spf13/cast"

	"uniswap-transaction/pkg/auth"
	"uniswap-transaction/pkg/metadata"
	"uniswap-transaction/pkg/notify"
	"uniswap-transaction/pkg/utils"
	"uniswap-transaction/pkg/utils/xtime"
)

// Option is metadata option.
type Option func(*options)

type options struct {
	envModes    []string
	ignoreCodes []int
}

// WithAlertEnvMode with env mode.
func WithAlertEnvMode(modes ...string) Option {
	return func(o *options) {
		o.envModes = modes
	}
}

// WithIgnoreCode with env mode.
func WithIgnoreCode(codes ...int) Option {
	return func(o *options) {
		o.ignoreCodes = codes
	}
}

type Manager struct {
	client   notify.Client
	template Template
	opts     options

	log *log.Helper
}

func NewAlertManager(client notify.Client, template Template, logger log.Logger, opts ...Option) *Manager {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	return &Manager{
		client:   client,
		template: template,
		opts:     o,
		log:      log.NewHelper(logger),
	}
}

func (a *Manager) SendServiceMessage(ctx context.Context, req interface{}, err error) {
	var (
		code      int32
		reason    string
		kind      string
		operation string
		userId    string
		username  string
		ts        = xtime.Now()
	)
	if err == nil || errors.Is(err, context.Canceled) {
		return
	}
	if info, ok := transport.FromServerContext(ctx); ok {
		kind = info.Kind().String()
		operation = info.Operation()
	}
	if claims, ok := auth.ClaimsFromContext(ctx); ok {
		userId = cast.ToString(claims.UserId)
		username = claims.Username
	}
	if se := errors.FromError(err); se != nil {
		code = se.Code
		reason = se.Reason
	}
	appInfo, _ := kratos.FromContext(ctx)
	envMode, _ := appInfo.Metadata()[metadata.EnvMode]
	if !utils.InSlice(envMode, a.opts.envModes) {
		return
	}
	if utils.InSliceInt(int(code), a.opts.ignoreCodes) {
		return
	}
	_, stack := extractError(err)
	m := map[string]interface{}{
		"environment": envMode,
		"ts":          ts.ToString(),
		"hostname":    appInfo.ID(),
		"service":     appInfo.Name(),
		"traceid":     tracing.TraceID()(ctx),
		"spanid":      tracing.SpanID()(ctx),
		"kind":        "server",
		"component":   kind,
		"operation":   operation,
		"userid":      userId,
		"username":    username,
		"args":        extractArgs(req),
		"code":        cast.ToString(code),
		"reason":      reason,
		"stack":       stack,
	}
	err = a.client.SendTemplateMessage(a.template.ServiceMessageContent(), m)
	if err != nil {
		a.log.WithContext(ctx).Errorf("alert error. err:%v", err)
	}
	return
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
