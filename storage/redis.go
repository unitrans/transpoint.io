package storage

import (
	"gopkg.in/redis.v3"
	"time"
)

type IRedisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value string, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
	TTL(key string) *redis.DurationCmd
}

type RedisDriver struct {
	Client *redis.Client
}

// TO remove public initiator
func Redis(address string) (*redis.Client) {

	var redisClient *redis.Client

	redisClient = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0, // use default DB
	})

	return redisClient
}

func NewRedisDriver(address string) (*RedisDriver) {
	return &RedisDriver{Client:Redis(address)}
}
