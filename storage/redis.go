package storage

import (
	"gopkg.in/redis.v3"
	"time"
	"errors"
)

type IRedisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value string, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
	TTL(key string) *redis.DurationCmd
}

type TranslationBag struct {
	Id           string `json:"id"`
	Source       string `json:"source"`
	Translations map[string]string `json:"translations"`
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

func (d *RedisDriver) GetLang(key, lang string) (bag TranslationBag, err error) {
	data, err := d.Client.HGetAllMap(key).Result()
	if nil != err {
		return
	}
	if _, exists := data[lang]; !exists {
		err = errors.New("Item not found")
		return
	}
	bag.Id = key
	bag.Source = data["source"]
	bag.Translations = map[string]string{lang:data[lang]}
	return
}

func (d *RedisDriver) GetAll(key string) (bag TranslationBag, err error) {
	data, err := d.Client.HGetAllMap(key).Result()
	if nil != err {
		return
	}

	bag.Id = key
	bag.Source = data["source"]
	delete(data, "source")
	bag.Translations = data
	return
}

func (d *RedisDriver) Save(key, source string, translations map[string]string) error {
	var transSlice []string
	for lang, trans := range translations {
		transSlice = append(transSlice, lang, trans)
	}
	return d.Client.HMSet(key, "source", source, transSlice...).Err()
}

func (d *RedisDriver) Delete(key string) (error) {
	return d.Client.Del(key).Err()
}

func (d *RedisDriver) DeleteLang(key, lang string) (error) {
	d.Client.HDel(key, lang).Err()
	data := d.Client.HKeys(key).Val()
	if 1 == len(data) {
		d.Client.Del(key)
	}
	return nil
}
