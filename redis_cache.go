package main

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(c *redis.Client, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: c,
		ttl:    ttl,
	}
}

func (c *RedisCache) Get(key int) (string, bool) {
	ctx := context.Background()
	keyStr := strconv.Itoa(key)
	val, err := c.client.Get(ctx, keyStr).Result()
	if err != nil {
		return val, false
	}
	return val, true
}

func (c *RedisCache) Set(key int, val string) error {
	ctx := context.Background()
	_, err := c.client.Set(ctx, strconv.Itoa(key), val, c.ttl).Result()
	return err
}

func (c *RedisCache) Remove(key int) error {
	ctx := context.Background()
	_, err := c.client.Del(ctx, strconv.Itoa(key)).Result()
	return err
}
