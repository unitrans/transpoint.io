// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package repository

import (
	"gopkg.in/redis.v3"
	"errors"
	"strings"
	"encoding/json"
)

func NewTranslationRepository(redis *redis.Client) *TranslationRepository {
	return &TranslationRepository{redis}
}

type TranslationRepository struct {
	client *redis.Client
}

type TranslationBag struct {
	Id           string `json:"id"`
	Source       string `json:"source"`
	Original     string `json:"original"`
	Meta     	 interface{} `json:"meta"`
	Translations map[string]string `json:"translations"`
}

func (t *TranslationBag) Langs() []string{
	var langs []string
	for k := range t.Translations{
		langs = append(langs, k)
	}
	return langs
}

func (t *TranslationBag) SetId(id string){
	split := strings.Split(id, "%")
	id = split[0]
	if len(split) > 1 {
		id = strings.Join(split[1:], "%")
	}
	t.Id = id
}



var (
	NotFoundError = errors.New("Item not found")
)

func (d *TranslationRepository) GetLang(key, lang string) (bag TranslationBag, err error) {
	data, err := d.client.HGetAllMap(key).Result()
	if nil != err {
		return
	}
	if _, exists := data[lang]; !exists {
		err = NotFoundError
		return
	}
	bag.SetId(key)
	bag.Source = data["source"]
	bag.Original = data["original"]
	bag.Translations = map[string]string{lang:data[lang]}
	return
}

func (d *TranslationRepository) GetAll(key string) (bag TranslationBag, err error) {
	data, err := d.client.HGetAllMap(key).Result()
	if nil != err {
		return
	}
	if 0 == len(data) {
		err = NotFoundError
		return
	}
	var metaDecoded interface{}
	if meta, ok := data["meta"]; ok {
		json.Unmarshal([]byte(meta), &metaDecoded)
	}

	bag.SetId(key)
	bag.Source = data["source"]
	bag.Original = data["original"]
	bag.Meta = metaDecoded
	delete(data, "meta")
	delete(data, "source")
	delete(data, "original")
	bag.Translations = data
	return
}

func (d *TranslationRepository) Save(key, source, original string, translations map[string]string, meta interface{}) error {
	var transSlice []string
	transSlice = append(transSlice, "original", original)
	metaEnc , _ := json.Marshal(meta)
	transSlice = append(transSlice, "meta", string(metaEnc))
	for lang, trans := range translations {
		transSlice = append(transSlice, lang, trans)
	}
	return d.client.HMSet(key, "source", source, transSlice...).Err()
}

func (d *TranslationRepository) Delete(key string) (error) {
	return d.client.Del(key).Err()
}

func (d *TranslationRepository) DeleteLang(key, lang string) (error) {
	d.client.HDel(key, lang).Err()
	data := d.client.HKeys(key).Val()
	if 2 == len(data) {
		d.client.Del(key)
	}
	return nil
}
