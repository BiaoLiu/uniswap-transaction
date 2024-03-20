package redis

import (
	"context"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	red "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"uniswap-transaction/pkg/utils/random"
)

const (
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
	lockCommand     = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
)

// A RedisLock is a redis lock.
type RedisLock struct {
	store   *Client
	seconds uint32
	key     string
	id      string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(store *Client, key string) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    random.Randn(randomLen),
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire(ctx context.Context) (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	resp, err := rl.store.Eval(ctx, lockCommand, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	}).Result()
	if err == red.Nil {
		return false, nil
	} else if err != nil {
		return false, errors.WithMessagef(err, "Error on acquiring lock for %s", rl.key)
	} else if resp == nil {
		return false, nil
	}
	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}
	// log.Errorf("Unknown reply when acquiring lock for %s: %v", rl.key, resp)
	return false, nil
}

// Release releases the lock.
func (rl *RedisLock) Release(ctx context.Context) (bool, error) {
	resp, err := rl.store.Eval(ctx, delCommand, []string{rl.key}, []string{rl.id}).Result()
	if err != nil {
		return false, err
	}
	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	return reply == 1, nil
}

// SetExpire sets the expiration.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}
