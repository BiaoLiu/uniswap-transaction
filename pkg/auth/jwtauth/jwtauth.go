package jwtauth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	jwtV4 "github.com/golang-jwt/jwt/v4"
	"golang.org/x/sync/singleflight"

	rdb "uniswap-transaction/pkg/cache/redis"
)

const (
	JwtTokenKey = "%sjwt_token:%d:%s"
)

var (
	ErrTokenInvalid = errors.Unauthorized("UNAUTHORIZED", "Token is invalid")
)

var (
	singleFlight = &singleflight.Group{}
)

type JwtAuth struct {
	SignKey []byte
	prefix  string
	store   *rdb.Client
}

func NewJwtAuth(prefix, signKey string, store *rdb.Client) *JwtAuth {
	return &JwtAuth{
		SignKey: []byte(signKey),
		prefix:  prefix,
		store:   store,
	}
}

// CreateToken 创建token
func (j *JwtAuth) CreateToken(claims *CustomClaims) (string, error) {
	jwtToken := jwtV4.NewWithClaims(jwtV4.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(j.SignKey)
	if err != nil {
		return "", err
	}
	key := fmt.Sprintf(JwtTokenKey, j.prefix, claims.UserId, token)
	prefixKey := fmt.Sprintf(JwtTokenKey, j.prefix, claims.UserId, "*")
	keys, err := j.store.Keys(context.Background(), prefixKey).Result()
	if err != nil {
		return "", err
	}
	if len(keys) > 0 {
		if err = j.store.Del(context.Background(), keys...).Err(); err != nil {
			return "", err
		}
	}
	if err = j.store.Set(context.Background(), key, &claims.BaseClaims, claims.ExpiresAt.Sub(time.Now())).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func (j *JwtAuth) DeleteToken(claims *CustomClaims, token string) error {
	key := fmt.Sprintf(JwtTokenKey, j.prefix, claims.UserId, token)
	_, err := j.store.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (j *JwtAuth) ValidateToken(claims *CustomClaims, token string) (bool, error) {
	key := fmt.Sprintf(JwtTokenKey, j.prefix, claims.UserId, token)
	count, err := j.store.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// CreateTokenByOldToken 旧token换新token，使用归并回源避免并发问题
func (j *JwtAuth) CreateTokenByOldToken(oldToken string, claims *CustomClaims) (string, error) {
	v, err, _ := singleFlight.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	return v.(string), err
}
