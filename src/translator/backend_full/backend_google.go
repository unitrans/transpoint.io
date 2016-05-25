// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package backend_full


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

func NewGoogleTranslator(client *http.Client, key string) IBackendFull {
	return &GoogleTranslator{client:client,key:key}
}


func (t *GoogleTranslator) TranslateFull(text string, language string) (IBackendFullResponse){
	data := &GoogleResponse{}
	data.Lang = language
	req, _ := http.NewRequest("GET", GT_URL+"?"+ t.getQueryString(text, language), nil)

	resp, err := t.client.Do(req)
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
	form.Add("format", "text")
	form.Add("prettyprint", "false")
	return form.Encode()
}



func(r *GoogleTranslator) GetName() string{
	return "google"
}

type GoogleResponse struct {
	Lang string
	Data *struct {
		Translations []struct{
			Text string `json:"translatedText"`
			Source string  `json:"detectedSourceLanguage"`
		} `json:"translations"`
	} `json:"data"`
	Error *struct{
		Code int `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func(r *GoogleResponse) GetText() string{
	if r.Error != nil {
		return ""
	}
	return r.Data.Translations[0].Text
}

func(r *GoogleResponse) GetSource() string{
	if r.Error != nil {
		return ""
	}
	return r.Data.Translations[0].Source
}

func(r *GoogleResponse) GetLang() string{
	if r.Error != nil {
		return ""
	}
	return r.Lang
}

func(r *GoogleResponse) GetName() string{
	return "google"
}
func(r *GoogleResponse) IsOk() bool {
	return r.Error == nil
}

