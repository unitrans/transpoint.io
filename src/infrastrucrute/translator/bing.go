// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translator
import (
	"net/http"
"log"
	"encoding/json"
	"net/url"
)

const BG_URL = "https://api.datamarket.azure.com/Bing/MicrosoftTranslator/v1/Translate"

type BingTranslator struct {
	client *http.Client
	key string
}

func NewBingTranslator(key string) ITranslateBackend {
	return &BingTranslator{initClient(), key}
}

func (t *BingTranslator) TranslateOne(text string, language string) (IBackendResponse){
	data := &BingResponse{}
	data.Lang = language
	req, _ := http.NewRequest("GET", BG_URL+"?"+ t.getQueryString(text, language), nil)
	req.SetBasicAuth(t.key, t.key)

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

func (t *BingTranslator) getQueryString(text, lang string)string{
	form := url.Values{}
	form.Add("To", "'"+lang+"'")
	form.Add("Text", "'"+text+"'")
	form.Add("$format", "json")
	return form.Encode()
}



func(r *BingTranslator) GetName() string{
	return "bing"
}

type BingResponse struct {
	Lang string
	Data struct {
		Results []struct{
			Text string `json:"Text"`
		} `json:"results"`
	} `json:"d"`
}

func(r *BingResponse) GetText() string{
	return r.Data.Results[0].Text
}

func(r *BingResponse) GetSource() string{
	return "na"
}

func(r *BingResponse) GetLang() string{
	return r.Lang
}

func(r *BingResponse) GetName() string{
	return "bing"
}