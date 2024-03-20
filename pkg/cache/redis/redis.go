package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config client settings.
type Config struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Client struct {
	*redis.Client
}

// NewRedisCmd new redis cmd client.
func NewRedisCmd(c *Config) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		PoolSize:     c.PoolSize,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
		DialTimeout:  time.Second * 2,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		panic(fmt.Sprintf("redis connection error. %s", err.Error()))
	}
	return &Client{Client: client}
}

func (c *Client) InPipeline(ctx context.Context, db int, fn func(pipe redis.Pipeliner)) ([]redis.Cmder, error) {
	pipe := c.Select(ctx, db)
	fn(pipe)
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}
	if cmds[0].Err() != nil {
		return nil, cmds[0].Err()
	}
	if len(cmds) > 1 {
		return cmds[1:], err
	}
	return cmds, err
}

func (c *Client) Select(ctx context.Context, db int) redis.Pipeliner {
	p := c.TxPipeline()
	p.Do(ctx, "select", db)
	return p
}
