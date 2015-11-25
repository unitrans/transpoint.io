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
)

// Translator interface
type Translator interface {
	Translate(text string, languages []string) *TranslationContainer
}

// TranslationBag hashmap
type TranslationBag map[string]string

// TranslationContainer struct
type TranslationContainer struct {
	Translations TranslationBag
	Source       string
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
}

type TranslateAdapter struct {
	client *http.Client
	translateBackend ITranslateBackend
}

func NewTranslateAdapter(back ITranslateBackend) *TranslateAdapter {
	return &TranslateAdapter{client:initClient(), translateBackend:back}
}

func (t *TranslateAdapter) Translate(text string, langs []string) *TranslationContainer {
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
	processor := NewEmojiProcessor()
	text = processor.Process(text)

	responseChan := make(chan IBackendResponse, len(languages))

	go t.doRequests(text, languages, responseChan)
	for resp := range responseChan {
		log.Println(resp)
		container.Translations[resp.GetLang()] = processor.Restore(resp.GetText())
		container.Source = resp.GetSource()
	}
	return container
}

func (t *TranslateAdapter) doRequests(text string, languages []string, c chan<- IBackendResponse){
	wg := &sync.WaitGroup{}
	for _, v := range languages {
		wg.Add(1)
		go func(text, lang string){
			defer wg.Done()
			resp := t.translateBackend.TranslateOne(text, lang)
			log.Println(lang, text, resp)
			c <- resp
		}(text, v)
	}
	wg.Wait()
	close(c)
}