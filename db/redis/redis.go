package redis

import (
	"context"
	"learn/settings"

	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var Client *RedisClient

func (r *RedisClient) Set(key string, value interface{}) {
	r.client.Set(r.ctx, key, value, 0)
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func MakeClient() {
	Client = &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     settings.Settings.Redis.Url(),
			Password: "", // no password set
			DB:       settings.Settings.Redis.DB,
		}),
		ctx: context.Background(),
	}
}
