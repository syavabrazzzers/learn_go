package redis

import (
	"context"
	"encoding/json"
	"learn/settings"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

var Client *RedisClient

func (r *RedisClient) Set(key string, value interface{}) {
	result := r.client.Set(r.ctx, key, value, 0)
	if result.Err() != nil {
		panic(result.Err())
	}
}

func (r *RedisClient) SetJson(key string, value interface{}, expiration int) {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	set_err := r.client.Set(r.ctx, key, data, time.Duration(time.Duration(expiration).Minutes())).Err()
	if set_err != nil {
		panic(set_err)
	}
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) GetJson(key string) (map[string]string, error) {
	var result map[string]string
	jsonStr, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		log.Fatalf("Error getting key: %v", err)
	}
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		log.Fatalf("JSON Unmarshal error: %v", err)
	}
	return result, nil
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
