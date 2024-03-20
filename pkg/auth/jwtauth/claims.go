package jwtauth

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	ExpiresTime time.Duration
	BufferTime  time.Duration
	Issuer      string
}

type CustomClaims struct {
	BaseClaims           `json:"baseClaims"`
	BufferTime           int64 `json:"bufferTime"`
	jwt.RegisteredClaims `json:"registeredClaims"`
}

func (c *CustomClaims) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(c)
	return
}

type BaseClaims struct {
	UserId      int64  `json:"userId"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AuthorityId string `json:"authorityId"`
}

func (c *BaseClaims) MarshalBinary() (data []byte, err error) {
	data, err = json.Marshal(c)
	return
}

func NewCustomClaims(c *Config, baseClaims BaseClaims) *CustomClaims {
	claims := &CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(c.BufferTime.Seconds()), // 缓冲时间
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Second * 1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.ExpiresTime)),       // 过期时间
			Issuer:    c.Issuer,                                                // 签名的发行者
		},
	}
	return claims
}

func (c *CustomClaims) CheckRenew() bool {
	if c.ExpiresAt.Time.Sub(time.Now()).Seconds() > float64(c.BufferTime) {
		return false
	}
	return true
}
