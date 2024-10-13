package database

import (
	"github.com/guirialli/rater_limit/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	r := config.LoadRedisConfig()
	return redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.Db,
	})
}
