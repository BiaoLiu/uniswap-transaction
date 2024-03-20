package redis

import (
	"testing"
)

func TestNewRedis(t *testing.T) {
	c := &Config{
		Addr:     "127.0.0.1:3306",
		Password: "123456",
		DB:       0,
	}
	rdb := NewRedisCmd(c)
	rdb.Close()
}
