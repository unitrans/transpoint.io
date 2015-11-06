// Copyright 2015 Yury Kozyrev. All rights reserved.
// Proprietary license.

//Package storage infrastructure/redis
package storage
import (
	"gopkg.in/redis.v3"
	"sync"
)
var (
	redisClient *redis.Client
	once = &sync.Once{}
)

// RedisClient singleton
func RedisClient(address, password string) (*redis.Client) {
	once.Do(func(){
		redisClient = redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       0,
		})
	})

	return redisClient
}