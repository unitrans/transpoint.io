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
//	"sync"
	"sync"
)

const YT_KEY = "trnsl.1.1.20150414T114844Z.a033fd4e2f954a35.ae3277873ab19d5355dbced6d668dc71b3011865"
const YT_URL = "https://translate.yandex.net/api/v1.5/tr.json/translate"

type YandexTranslator struct {
	client *http.Client
}

func NewYandexTranslator() *YandexTranslator {
	return &YandexTranslator{client:initClient()}
}

func (t *YandexTranslator) Translate(text string, langs []string) (*TranslationContainer) {
	container := &TranslationContainer{
		Translations:TranslationBag{},
	}

	set := make(map[string]bool)
	for _, val := range langs {
		set[val] = true
	}
	var languages []string
	for lang, _ := range set{
		languages = append(languages, lang)
	}

	responseChan := make(chan *YandexResponse, len(languages))

	go t.doRequests(text, languages, responseChan)
	for resp := range responseChan {
		log.Println(resp)
		container.Translations[resp.GetLang()] = resp.GetText()
		container.Source = resp.GetSource()
	}
	return container
}

func (t *YandexTranslator) doRequests(text string, languages []string, c chan *YandexResponse){
	wg := &sync.WaitGroup{}
	for _, v := range languages {
		wg.Add(1)
		go func(text, lang string){
			defer wg.Done()
			resp := t.TranslateOne(text, lang)
			log.Println(lang, text, resp)
			c <- resp
		}(text, v)
	}
	wg.Wait()
	close(c)
}

func (t *YandexTranslator) TranslateOne(text string, language string) (data *YandexResponse){
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
	return
}

func (t *YandexTranslator) getQueryString(text, lang string)string{
	form := url.Values{}
	form.Add("key", YT_KEY)
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