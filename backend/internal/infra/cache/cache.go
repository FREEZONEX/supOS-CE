package cache

import (
	"context"
	"errors"
	"time"

	"backend/internal/config"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

var ErrNotFound = errors.New("cache not found")

func Open(c config.RedisConf) *Client {
	return &Client{rdb: redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})}
}

func (c *Client) Close() error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.rdb == nil {
		return errors.New("redis not initialized")
	}
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Publish(ctx context.Context, channel string, payload any) error {
	if c == nil || c.rdb == nil {
		return errors.New("redis not initialized")
	}
	return c.rdb.Publish(ctx, channel, payload).Err()
}

func (c *Client) GetString(ctx context.Context, key string) (string, error) {
	if c == nil || c.rdb == nil {
		return "", errors.New("redis not initialized")
	}
	value, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	return value, err
}

func (c *Client) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	if c == nil || c.rdb == nil {
		return errors.New("redis not initialized")
	}
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if c == nil || c.rdb == nil {
		return errors.New("redis not initialized")
	}
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Client) EvalInt64(ctx context.Context, script string, keys []string, args ...any) (int64, error) {
	if c == nil || c.rdb == nil {
		return 0, errors.New("redis not initialized")
	}
	return c.rdb.Eval(ctx, script, keys, args...).Int64()
}
