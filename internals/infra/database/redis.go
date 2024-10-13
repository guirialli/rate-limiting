package database

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/guirialli/rater_limit/config"
	"github.com/redis/go-redis/v9"
)

type RedisDatabase[T any] struct {
	rdb *redis.Client
}

func NewRedisClient[T any](config config.Redis) *RedisDatabase[T] {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.Db,
	})
	return &RedisDatabase[T]{client}
}

func (r *RedisDatabase[T]) Get(ctx context.Context, key string) (*T, bool) {
	val, err := r.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	} else if err != nil {
		panic(err)
	}

	var result T
	if err = json.Unmarshal([]byte(val), &result); err != nil {
		panic(err)
	}
	return &result, true
}

func (r *RedisDatabase[T]) Set(ctx context.Context, key string, value *T) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}

	r.rdb.Set(ctx, key, val, 0)
	return nil
}
