// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package backend_particular

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const YA_DICT_URL = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup"

//https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=dict.1.1.xxx&lang=de-ru&text=fertig

type YandexDict struct {
	client *http.Client
	key    string
}

func NewYandexDict(c *http.Client, key string) IBackendParticular {
	return &YandexDict{client: c, key: key}
}

func (t *YandexDict) TranslateWord(text string, language, to string) IBackendParticularResponse {

	data := &YandexDictResponse{}
	data.Lang = language

	reqUrl := YA_DICT_URL + "?" + t.getQueryString(text, language, to)
	req, _ := http.NewRequest("GET", reqUrl, nil)
	data.Url = reqUrl

	resp, err := t.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	reader := ioutil.NopCloser(resp.Body)

	// log.Println(str)
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		log.Println("error decode", err)
	}

	return data
}

func (t *YandexDict) GetName() string {
	return "yandex_dict"
}

func (t *YandexDict) getQueryString(text, from, to string) string {
	form := url.Values{}
	form.Add("key", t.key)
	form.Add("lang", from+"-"+to)
	form.Add("text", text)
	return form.Encode()
}

type YandexDictResponse struct {
	Lang string
	Url  string
	Def  []*YandexDictResponseItem `json:"def"`
}
type YandexDictResponseItem struct {
	Text string                    `json:"text"`
	Num  string                    `json:"pos"`
	Def  []*YandexDictResponseItem `json:"def"`
	Tr   []*YandexDictResponseItem `json:"tr"`
	Syn  []*YandexDictResponseItem `json:"syn"`
	Mean []*YandexDictResponseItem `json:"mean"`
}

func (t *YandexDictResponse) GetUrl() string {
	return t.Url
}

func (t *YandexDictResponse) GetMeanings() []IParticularMeaning {
	meanings := []IParticularMeaning{}
	for _, v := range t.Def {
		for _, tr := range v.Tr {
			meaning := &Meaning{Text:tr.Text}
			meanings = append(meanings, meaning)
		}

	}
	return meanings
}
