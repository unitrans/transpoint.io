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
	"github.com/urakozz/transpoint.io/src/infrastrucrute/translator/particular"
	"regexp"
)

// Translator interface
type Translator interface {
	Translate(text string, languages []string) *TranslationContainer
	AddParticular(p particular.IParticularBackend)
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
	translateParticular []particular.IParticularBackend
//	middleware []translation_middleware.ITranslationMiddleware
	landChan chan string

}

func NewTranslateAdapter(back []ITranslateBackend) *TranslateAdapter {
	return &TranslateAdapter{client:initClient(), translateBackend:back, landChan:make(chan string, 1)}
}
func (t *TranslateAdapter) AddParticular(p particular.IParticularBackend){
	p.SetClient(t.client)
	t.translateParticular = append(t.translateParticular, p)
}
func (t *TranslateAdapter) Translate(text string, langs []string) *TranslationContainer {
	container := t.newContainer()
	container.Original = text

	langs = array_uniq(langs)
	processor := NewEmojiProcessor()
	text = processor.Process(text)

	responseChan := make(chan *RawTransData, len(langs))
	responseChanParticular := make(chan *RawParticularData, len(langs))

	go t.doRequestsTranslators(text, langs, responseChan)
	go func(){
		lang := <-t.landChan
		t.doRequestsParticular(text, lang, langs, responseChanParticular)
	}()
	for resp := range responseChan {
		log.Println(resp)
		container.RawTranslations[resp.Name][resp.Lang] = resp.Translation
		if "google" == resp.Name {
			container.Source = resp.Source
			t.landChan <- resp.Source
		}
		resp.Translation = processor.Restore(resp.Translation)
		container.RawTransData = append(container.RawTransData, resp)
	}
	container.Translations = container.RawTranslations["google"]

	for resp := range responseChanParticular {
		container.RawParticularData = append(container.RawParticularData, resp)
	}
	return container
}

func (t *TranslateAdapter) doRequestsTranslators(text string, languages []string, c chan<- *RawTransData){
	emFix := regexp.MustCompile(`\_\s?\$\s?(\d+)\_+`)
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
					Translation:emFix.ReplaceAllString(resp.GetText(), `_\$$1_`),
					Time: time.Since(t) / time.Millisecond,
				}
				c <- raw
			}(text, v, back)
		}

	}
	wg.Wait()
	close(c)
}

func (t *TranslateAdapter) doRequestsParticular(text, source string, languages []string, c chan<- *RawParticularData){

	bag := regexp.MustCompile(`([\w\x7f-\xff])+`).FindAllString(text, -1)

	wg := &sync.WaitGroup{}
	for _, lang := range languages {
		for _, back := range t.translateParticular{
			wg.Add(1)
			raw := &RawParticularData{
				Original: text,
				Source:source,
				Lang:lang,
				Name: back.GetName(),
				ParticularBag: make([]*ParticularItem, len(bag)),
			}
			go func(raw *RawParticularData, backend particular.IParticularBackend){
				defer wg.Done()
				t := time.Now()

				partChan := make(chan *ParticularItem, len(bag))
				wgSingle := &sync.WaitGroup{}
				for index, word := range bag {
					wgSingle.Add(1)
					item := &ParticularItem{Order: index, Original: word}
					go func(item *ParticularItem,raw *RawParticularData, backend particular.IParticularBackend, c chan<- *ParticularItem){
						defer wgSingle.Done()
						t1 := time.Now()
						resp := backend.TranslateOne(item.Original, raw.Source, raw.Lang)
						item.Translations = resp.GetMeanings()
						item.Url = resp.GetUrl()
						item.Time = time.Since(t1) / time.Millisecond
						c <- item
					}(item, raw, backend, partChan)
				}
				wgSingle.Wait()
				close(partChan)

				for item := range partChan {
					raw.ParticularBag[item.Order] = item
				}

				raw.Time = time.Since(t) / time.Millisecond

				c <- raw
			}(raw, back)
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