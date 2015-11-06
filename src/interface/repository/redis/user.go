// Copyright 2015 Yury Kozyrev. All rights reserved.
// Proprietary license.
package repository
import "gopkg.in/redis.v3"

func NewUserRepository(redis *redis.Client) *UserRepository {
	return &UserRepository{redis}
}

type UserRepository struct {
	client *redis.Client
}

func (r *UserRepository) GetSecretByKey(key string) string {
	return r.client.HGet("keys", key).Val()
}
func (r *UserRepository) GetAllSecretsByKeys(key... string) ([]interface{}, error) {
	return r.client.HMGet("keys", key...).Result()
}
func (r *UserRepository) DeleteSecretByKey(key string) error {
	return r.client.HDel("keys", key).Err()
}
func (r *UserRepository) SaveSecretByKey(key, secret string) error {
	return r.client.HSet("keys", key, secret).Err()
}

func (r *UserRepository) GetUserById(id string) (string, error) {
	return r.client.HGet("user", id).Result()
}
func (r *UserRepository) SaveUserById(id, data string) error {
	return r.client.HSet("user", id, data).Err()
}

