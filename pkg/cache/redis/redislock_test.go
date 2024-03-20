package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"kratos-sharp/pkg/utils"
)

func newRedis() *Client {
	c := &Config{
		Addr:     "127.0.0.1:6379",
		Password: "findcbi@com",
		DB:       0,
	}
	rdb := NewRedisCmd(c)
	return rdb
}

func TestRedisLock(t *testing.T) {
	testFn := func(ctx context.Context) func(client *Client) {
		return func(client *Client) {
			key := utils.Rand()
			firstLock := NewRedisLock(client, key)
			firstLock.SetExpire(60)
			firstAcquire, err := firstLock.Acquire(context.Background())
			assert.Nil(t, err)
			assert.True(t, firstAcquire)

			secondLock := NewRedisLock(client, key)
			secondLock.SetExpire(5)
			againAcquire, err := secondLock.Acquire(context.Background())
			assert.Nil(t, err)
			assert.False(t, againAcquire)

			release, err := firstLock.Release(context.Background())
			assert.Nil(t, err)
			assert.True(t, release)

			endAcquire, err := secondLock.Acquire(context.Background())
			assert.Nil(t, err)
			assert.True(t, endAcquire)
		}
	}

	f := testFn(context.Background())
	f(newRedis())
}

func TestRedisLock_Expired(t *testing.T) {
	client := newRedis()
	testFn := func(client *Client) {
		key := utils.Rand()
		redisLock := NewRedisLock(client, key)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := redisLock.Acquire(ctx)
		assert.NotNil(t, err)
	}
	testFn(client)

	testFn = func(client *Client) {
		key := utils.Rand()
		redisLock := NewRedisLock(client, key)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := redisLock.Release(ctx)
		assert.NotNil(t, err)
	}
	testFn(client)
}
