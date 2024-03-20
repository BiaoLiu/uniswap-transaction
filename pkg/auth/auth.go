package auth

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"uniswap-transaction/pkg/auth/jwtauth"
	"uniswap-transaction/pkg/middleware/casbin"
	"uniswap-transaction/pkg/middleware/jwt"
)

type SecurityUser struct {
	Path        string
	Method      string
	AuthorityId string
	Domain      string
}

func NewSecurityUser() casbin.SecurityUser {
	return &SecurityUser{}
}

func (su *SecurityUser) ParseFromContext(ctx context.Context) error {
	if claims, ok := ClaimsFromContext(ctx); ok {
		su.AuthorityId = claims.AuthorityId
		//su.Domain = claims.(jwtV4.MapClaims)[Domain].(string)
	} else {
		return errors.New("jwt claim missing")
	}

	if tr, ok := transport.FromServerContext(ctx); ok {
		httpTr := tr.(*http.Transport)
		su.Path = httpTr.PathTemplate()
		su.Method = httpTr.Request().Method
	} else {
		return errors.New("jwt claim missing")
	}
	return nil
}

func (su *SecurityUser) GetSubject() string {
	return su.AuthorityId
}

func (su *SecurityUser) GetObject() string {
	return su.Path
}

func (su *SecurityUser) GetAction() string {
	return su.Method
}

func (su *SecurityUser) GetDomain() string {
	return su.Domain
}

// ClaimsFromContext extract auth info from context
func ClaimsFromContext(ctx context.Context) (customClaims *jwtauth.CustomClaims, ok bool) {
	if claims, ok := jwt.ClaimsFromContext(ctx); ok {
		customClaims, ok = claims.(*jwtauth.CustomClaims)
		return customClaims, ok
	}
	return nil, false
}

// TokenFromContext extract auth info from context
func TokenFromContext(ctx context.Context) (token string, ok bool) {
	return jwt.TokenFromContext(ctx)
}
