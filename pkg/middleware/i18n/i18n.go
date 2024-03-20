package i18n

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"golang.org/x/text/language"
)

type localizeKey struct{}

func I18N(localize func(lang language.Tag)) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				accept := tr.RequestHeader().Get("accept-language")
				println(accept)
				lang, err := language.Parse(accept)
				if err != nil {
					lang = language.English
				}
				if localize != nil {
					localize(lang)
				}
				ctx = context.WithValue(ctx, localizeKey{}, lang)
			}
			return handler(ctx, req)
		}
	}
}
