package redisinfra

import (
	"context"
	"fmt"
	"time"

	"github.com/datpham2001/be-winsku/internal/infrastructure/configinfra"
	"github.com/redis/go-redis/v9"
)

const DEFAULT_TTL = 15 * time.Minute

type Client struct {
	*redis.Client
}

func NewConnection(cfg configinfra.RedisConfig) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	fmt.Println("connect to redis successfully")
	return &Client{client}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) HealthCheck() error {
	return c.Client.Ping(context.Background()).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.Client.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}
