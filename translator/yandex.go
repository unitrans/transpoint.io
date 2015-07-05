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

type YandexTranslator struct {
	client *http.Client
}

func NewYandexTranslator() *YandexTranslator {
	return &YandexTranslator{client:initClient()}
}

func (t *YandexTranslator) Translate(text string, languages []string) *TranslationContainer {
	form := url.Values{}
	form.Add("key", "trnsl.1.1.20150414T114844Z.a033fd4e2f954a35.ae3277873ab19d5355dbced6d668dc71b3011865")
	form.Add("lang", languages[0])
	form.Add("text", text)
	req, _ := http.NewRequest("POST", "https://translate.yandex.net/api/v1.5/tr.json/translate", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	var data *YandexResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	log.Println("response obj:", data)
	return &TranslationContainer{
		Source: data.GetSource(),
		Bag:TranslationBag{data.GetLang():data.GetText()},
	}
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