package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lejenome/lro/pkg/config"
)

var CTX context.Context

type Redis struct {
	db *redis.Client
}

func New(config *config.RedisConfig) *Redis {
	DB := redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})
	CTX = context.Background()
	return &Redis{
		db: DB,
	}
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.db.Set(CTX, key, value, expiration)
}
func (r *Redis) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.db.SetNX(CTX, key, value, expiration)
}
func (r *Redis) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.db.SetXX(CTX, key, value, expiration)
}
func (r *Redis) Exists(keys ...string) *redis.IntCmd {
	return r.db.Exists(CTX, keys...)
}
