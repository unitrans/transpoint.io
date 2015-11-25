// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package translator
import (

	"net/url"
	"net/http"
	"bytes"
	"log"
	"encoding/json"
	"strings"
)

const YT_URL = "https://translate.yandex.net/api/v1.5/tr.json/translate"

type YandexTranslator struct {
	client *http.Client
	key string
}

func NewYandexTranslator(key string) ITranslateBackend {
	return &YandexTranslator{client:initClient(), key:key}
}

func (t *YandexTranslator) TranslateOne(text string, language string) (IBackendResponse){
	data := &YandexResponse{}

	req, _ := http.NewRequest("POST", YT_URL, bytes.NewBufferString(t.getQueryString(text, language)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
	}
	return data
}

func (t *YandexTranslator) getQueryString(text, lang string)string{
	form := url.Values{}
	form.Add("key", t.key)
	form.Add("lang", lang)
	form.Add("text", text)
	return form.Encode()
}

type YandexResponse struct {
	Code int `json:"code"`
	Lang string `json:"lang"`
	Text []string `json:"text"`
}

func(r *YandexResponse) GetText() string{
	return r.Text[0]
}
func(r *YandexResponse) GetLang() string{
	return strings.Split(r.Lang, "-")[1]
}
func(r *YandexResponse) GetSource() string{
	return strings.Split(r.Lang, "-")[0]
}