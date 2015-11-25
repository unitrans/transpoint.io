// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translator
import (
	"net/http"
"log"
	"encoding/json"
	"net/url"
)

const GT_URL = "https://www.googleapis.com/language/translate/v2"

type GoogleTranslator struct {
	client *http.Client
	key string
}

func NewGoogleTranslator(key string) ITranslateBackend {
	return &GoogleTranslator{client:initClient(), key:key}
}

func (t *GoogleTranslator) TranslateOne(text string, language string) (IBackendResponse){
	data := &GoogleResponse{}
	data.Lang = language
	req, _ := http.NewRequest("GET", GT_URL+"?"+ t.getQueryString(text, language), nil)

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

func (t *GoogleTranslator) getQueryString(text, lang string)string{
	form := url.Values{}
	form.Add("key", t.key)
	form.Add("target", lang)
	form.Add("q", text)
	return form.Encode()
}

type GoogleResponse struct {
	Lang string
	Data struct {
		Translations []struct{
			Text string `json:"translatedText"`
			Source string  `json:"detectedSourceLanguage"`
		} `json:"translations"`
	} `json:"data"`
}

func(r *GoogleResponse) GetText() string{
	return r.Data.Translations[0].Text
}

func(r *GoogleResponse) GetSource() string{
	return r.Data.Translations[0].Source
}

func(r *GoogleResponse) GetLang() string{
	return r.Lang
}