// Copyright 2015 Yury Kozyrev. All rights reserved.
// Proprietary license.

// Package translator
package translator
import (
	"net/http"
	"time"
	"github.com/facebookgo/httpcontrol"
	"sync"
	"log"
//	"github.com/urakozz/transpoint.io/src/infrastrucrute/translator/translation_middleware"
)

// Translator interface
type Translator interface {
	Translate(text string, languages []string) *TranslationContainer
}

var (
	client *http.Client
	once = &sync.Once{}
)

func initClient() (*http.Client) {
	once.Do(func() {
		transport := &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries: 3,
			MaxIdleConnsPerHost: 10,
		}

		client = &http.Client{
			Transport: transport,
		}
	})

	return client
}

type IBackendResponse interface {
	GetText() string
	GetSource() string
	GetLang() string
}

type ITranslateBackend interface {
	TranslateOne(text string, language string) (data IBackendResponse)
	GetName() string
}

type TranslateAdapter struct {
	client *http.Client
	translateBackend []ITranslateBackend
//	middleware []translation_middleware.ITranslationMiddleware

}

func NewTranslateAdapter(back []ITranslateBackend) *TranslateAdapter {
	return &TranslateAdapter{client:initClient(), translateBackend:back}
}

func (t *TranslateAdapter) Translate(text string, langs []string) *TranslationContainer {
	container := t.newContainer()
	container.Original = text

	langs = array_uniq(langs)
//	processor := NewEmojiProcessor()
//	text = processor.Process(text)

	responseChan := make(chan *RawTransData, len(langs))

	go t.doRequests(text, langs, responseChan)
	for resp := range responseChan {
		log.Println(resp)
		container.RawTranslations[resp.Name][resp.Lang] = resp.Translation
		if "google" == resp.Name {
			container.Source = resp.Source
		}
		container.RawTransData = append(container.RawTransData, resp)
	}
	container.Translations = container.RawTranslations["google"]
	return container
}

func (t *TranslateAdapter) doRequests(text string, languages []string, c chan<- *RawTransData){
	wg := &sync.WaitGroup{}
	for _, v := range languages {
		for _, back := range t.translateBackend{
			wg.Add(1)
			go func(text, lang string, backend ITranslateBackend){
				defer wg.Done()
				t := time.Now()
				resp := backend.TranslateOne(text, lang)

				raw := &RawTransData{
					Source:resp.GetSource(),
					Lang:resp.GetLang(),
					Name:backend.GetName(),
					Translation:resp.GetText(),
					Time: time.Since(t) / time.Millisecond,
				}
				c <- raw
			}(text, v, back)
		}

	}
	wg.Wait()
	close(c)
}

func (t *TranslateAdapter) newContainer() *TranslationContainer{
	raw := make(map[string]TranslationBag, len(t.translateBackend))
	for _, back := range t.translateBackend {
		raw[back.GetName()] = TranslationBag{}
	}
	return &TranslationContainer{
		Translations:TranslationBag{},
		RawTranslations:raw,
	}
}

func array_uniq(langs []string) []string {
	set := make(map[string]struct{})
	for _, val := range langs {
		set[val] = struct{}{}
	}
	var languages []string
	for lang, _ := range set{
		languages = append(languages, lang)
	}
	return languages
}