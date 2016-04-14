// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package translator

import (
	"github.com/urakozz/transpoint.io/src/translator/backend_particular"
	"github.com/urakozz/transpoint.io/src/translator/backend_full"
	"sync"
	"time"
	"github.com/urakozz/transpoint.io/src/components"
)

type Translator interface {
	Translate(text string, languages []string) *TranslationContainer
	//AddParticular(p particular.IParticularBackend)
}

type TranslateAdapter struct {
	translateBackend []backend_full.IBackendFull
	translateParticular []backend_particular.IBackendParticular
	markov components.IChain
	landChan chan string

}

func NewTranslateAdapter(back []backend_full.IBackendFull, markov components.IChain) *TranslateAdapter {
	return &TranslateAdapter{translateBackend:back, landChan:make(chan string, 1), markov:markov}
}

func (t *TranslateAdapter) Translate(text string, langs []string) *TranslationContainer {
	langs = array_uniq(langs)

	container := &TranslationContainer{
		Translations:TranslationBag{},
		RawTransData: make(map[string]map[string]*RawTranslationData),
	}
	container.Original = text

	//fill markov model
	if len(text) > 100{
		t.markov.Add(text)
	}

	responseChan := make(chan *RawTranslationData, len(langs))

	go t.doRequestsTranslators(text, langs, responseChan)

	for resp := range responseChan {
		if _, ok := container.RawTransData[resp.Lang]; !ok {
			container.RawTransData[resp.Lang] = make(map[string]*RawTranslationData)
		}
		resp.Score = t.markov.Occurrences(resp.Translation)
		container.RawTransData[resp.Lang][resp.Name] = resp
	}


	for lang, _ := range container.RawTransData{
		container.RawTransData[lang]["uni"] = &RawTranslationData{}
	}
	// set google by default
	for lang, details := range container.RawTransData{
		*container.RawTransData[lang]["uni"] = *details["google"]
	}
	// if google is empty (credential issues)
	for lang, details := range container.RawTransData{
		if details["google"].Translation == "" {
			*container.RawTransData[lang]["uni"] = *details["yandex"]
		}
	}
        // if yandex detects different lang, use it (also implicitly covers previous case)
	for lang, details := range container.RawTransData{
		if details["yandex"].Source != details["google"].Source {
			container.RawTransData[lang]["uni"].Translation = details["yandex"].Translation
			container.RawTransData[lang]["uni"].Source = details["yandex"].Source
			container.RawTransData[lang]["uni"].Lang = details["yandex"].Lang
		} else if details["yandex"].Score > details["google"].Score {
			container.RawTransData[lang]["uni"].Translation = details["yandex"].Translation
		}
	}
	//// fill markov model if confident source
	//for _, details := range container.RawTransData{
	//	if details["yandex"].Source == details["google"].Source {
	//		t.markov.Add()
	//	}
	//}

	//postprocessing
	for lang, details := range container.RawTransData{
		container.RawTransData[lang]["uni"].Name = "unitrans"
		container.RawTransData[lang]["uni"].Score = t.markov.Occurrences(details["uni"].Translation)

		container.Translations[lang] = details["uni"].Translation
		container.Source = details["uni"].Source
	}

	return container
}

func (t *TranslateAdapter) doRequestsTranslators(text string, languages []string, c chan<- *RawTranslationData){

	wg := &sync.WaitGroup{}
	for _, v := range languages {
		for _, back := range t.translateBackend{
			wg.Add(1)
			go func(text, lang string, backend backend_full.IBackendFull){
				defer wg.Done()
				t := time.Now()
				resp := backend.TranslateFull(text, lang)

				raw := &RawTranslationData{
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