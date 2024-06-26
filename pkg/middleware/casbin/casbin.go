package casbin

import (
	"context"

	stdcasbin "github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

type contextKey string

const (
	ModelContextKey        contextKey = "CasbinModel"
	PolicyContextKey       contextKey = "CasbinPolicy"
	EnforcerContextKey     contextKey = "CasbinEnforcer"
	SecurityUserContextKey contextKey = "CasbinSecurityUser"

	reason string = "FORBIDDEN"

	defaultRBACModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`
)

var (
	ErrSecurityUserCreatorMissing = errors.Forbidden(reason, "SecurityUserCreator is required")
	ErrEnforcerMissing            = errors.Forbidden(reason, "Enforcer is missing")
	ErrSecurityParseFailed        = errors.Forbidden(reason, "Security Info fault")
	ErrUnauthorized               = errors.Forbidden(reason, "Unauthorized Access")
)

type Option func(*options)

type options struct {
	enableDomain        bool
	securityUserCreator SecurityUserCreator
	model               model.Model
	policy              persist.Adapter
	enforcer            *stdcasbin.SyncedEnforcer
}

// WithDomainSupport  enable domain support
func WithDomainSupport() Option {
	return func(o *options) {
		o.enableDomain = true
	}
}

func WithSecurityUserCreator(securityUserCreator SecurityUserCreator) Option {
	return func(o *options) {
		o.securityUserCreator = securityUserCreator
	}
}

func WithCasbinModel(model model.Model) Option {
	return func(o *options) {
		o.model = model
	}
}

func WithCasbinPolicy(policy persist.Adapter) Option {
	return func(o *options) {
		o.policy = policy
	}
}

func WithCasbinEnforcer(enforcer *stdcasbin.SyncedEnforcer) Option {
	return func(o *options) {
		o.enforcer = enforcer
	}
}

// loadRbacModel 加载RBAC模型
func loadRbacModel() (model.Model, error) {
	return model.NewModelFromString(defaultRBACModel)
}

func Server(opts ...Option) middleware.Middleware {
	o := &options{
		securityUserCreator: nil,
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.model == nil {
		o.model, _ = loadRbacModel()
	}
	if o.enforcer == nil {
		o.enforcer, _ = stdcasbin.NewSyncedEnforcer(o.model, o.policy)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var (
				allowed bool
				err     error
			)

			if o.enforcer == nil {
				return nil, ErrEnforcerMissing
			}
			if o.securityUserCreator == nil {
				return nil, ErrSecurityUserCreatorMissing
			}

			securityUser := o.securityUserCreator()
			if err = securityUser.ParseFromContext(ctx); err != nil {
				return nil, ErrSecurityParseFailed
			}

			ctx = context.WithValue(ctx, SecurityUserContextKey, securityUser)
			if o.enableDomain {
				allowed, err = o.enforcer.Enforce(securityUser.GetSubject(), securityUser.GetDomain(), securityUser.GetObject(), securityUser.GetAction())
			} else {
				allowed, err = o.enforcer.Enforce(securityUser.GetSubject(), securityUser.GetObject(), securityUser.GetAction())
			}
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, ErrUnauthorized
			}
			return handler(ctx, req)
		}
	}
}

func Client(opts ...Option) middleware.Middleware {
	o := &options{
		securityUserCreator: nil,
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.model == nil {
		o.model, _ = loadRbacModel()
	}
	if o.enforcer == nil {
		o.enforcer, _ = stdcasbin.NewSyncedEnforcer(o.model, o.policy)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
	}
}

// SecurityUserFromContext extract SecurityUser from context
func SecurityUserFromContext(ctx context.Context) (SecurityUser, bool) {
	user, ok := ctx.Value(SecurityUserContextKey).(SecurityUser)
	return user, ok
}
