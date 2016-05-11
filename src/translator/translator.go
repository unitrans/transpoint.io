// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package translator

import (
	"github.com/unitrans/unitrans/src/components"
	"github.com/unitrans/unitrans/src/translator/backend_full"
	"github.com/unitrans/unitrans/src/translator/backend_particular"
	"github.com/unitrans/unitrans/src/translator/processing"
	"math"
	"strings"
	"sync"
	"time"
)

type Translator interface {
	Translate(text string, languages []string) *TranslationContainer
	//AddParticular(p particular.IParticularBackend)
}

type TranslateAdapter struct {
	translateBackend    []backend_full.IBackendFull
	translateParticular []backend_particular.IBackendParticular
	markov              components.IChain
	landChan            chan string
}

func NewTranslateAdapter(back []backend_full.IBackendFull, markov components.IChain) *TranslateAdapter {
	return &TranslateAdapter{translateBackend: back, landChan: make(chan string, 1), markov: markov}
}

func (t *TranslateAdapter) Translate(text string, langs []string) *TranslationContainer {
	langs = array_uniq(langs)

	container := t.createContainer(text, langs)

	//fill markov model
	if len(text) > 100 {
		t.markov.Add(text)
	}

	segments := processing.Segments.Split(text)
	container.RawSegmentsData = t.GetSegmentTranslations(segments, langs)
	t.CalculateScores(container, segments)
	t.AppendUniEngine(container, segments)
	t.CalculateRawTransData(container)

	for lang, details := range container.RawTransData {
		container.Translations[lang] = details["uni"].Translation
		container.Source = details["uni"].Source
	}

	return container
}

func (t *TranslateAdapter) CalculateScores(container *TranslationContainer, segments []*processing.Segment) {
	for i, seg := range segments {
		if seg.Type > processing.SegmentText {
			continue
		}
		for _, lang := range container.Langs {
			for _, back := range t.translateBackend {
				rawSegData := container.RawSegmentsData[i][lang][back.GetName()]
				rawSegData.Score = t.markov.Occurrences(rawSegData.Translation)
			}
		}
	}
}

func (t *TranslateAdapter) AppendUniEngine(container *TranslationContainer, segments []*processing.Segment) {

	for i, seg := range segments {
		for _, lang := range container.Langs {

			if seg.Type > processing.SegmentText {
				container.RawSegmentsData[i][lang]["uni"] = &RawTranslationData{Name: "uni", Translation: seg.Text}
			} else {
				container.RawSegmentsData[i][lang]["uni"] = t.ChooseSegment(container.RawSegmentsData[i][lang])
			}
		}
	}

}

func (t *TranslateAdapter) CalculateRawTransData(container *TranslationContainer) {

	for _, seg := range container.RawSegmentsData {
		for lang, langc := range seg {
			for back, backc := range langc {
				container.RawTransData[lang][back].Lang = lang
				container.RawTransData[lang][back].Translation += backc.Translation
				container.RawTransData[lang][back].Score += backc.Score
				if backc.Source != "" {
					currLang := container.RawTransData[lang][back].Source
					container.RawTransData[lang][back].Source = strings.Join(array_uniq(strings.Split(currLang+","+backc.Source, ",")), ",")
				}
			}
		}
	}
}

func (t *TranslateAdapter) ChooseSegment(engines map[string]*RawTranslationData) *RawTranslationData {

	result := &RawTranslationData{}
	// set yandex by default
	*result = *engines["google"]

	// otherwise google if yandex failed
	if engines["google"].Translation == "" {
		engines["google"].Score = float64(-math.MaxInt32)
	}

	// if yandex detects different lang, reduce google ranc
	if engines["yandex"].Source != engines["google"].Source {
		engines["google"].Score = float64(-math.MaxInt32)
	}
	maxScore := NewSegmentsSorter(engines).Max()
	*result = *maxScore
	result.Name = "uni"

	return result
}

func (t *TranslateAdapter) GetSegmentTranslations(textSegments []*processing.Segment, languages []string) []map[string]map[string]*RawTranslationData {

	translations := make([]map[string]map[string]*RawTranslationData, len(textSegments))
	for i, seg := range textSegments {
		translations[i] = make(map[string]map[string]*RawTranslationData, len(languages))
		for _, lang := range languages {
			translations[i][lang] = make(map[string]*RawTranslationData, len(t.translateBackend))
			for _, back := range t.translateBackend {
				translations[i][lang][back.GetName()] = &RawTranslationData{
					Name:     back.GetName(),
					Original: seg.Text,
				}
			}
		}

	}

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	//translationsSegments []*RawTranslationData
	for _, l := range languages {
		for _, b := range t.translateBackend {
			for i, s := range textSegments {
				wg.Add(1)
				go func(seg *processing.Segment, lang string, index int, backend backend_full.IBackendFull) {
					defer wg.Done()
					if seg.Type > processing.SegmentText {
						mu.Lock()
						translations[index][lang][backend.GetName()].Translation = seg.Text
						mu.Unlock()
					} else {
						t := time.Now()

						resp := backend.TranslateFull(seg.Text, lang)

						mu.Lock()
						translations[index][lang][backend.GetName()].Source = resp.GetSource()
						translations[index][lang][backend.GetName()].Lang = resp.GetLang()
						translations[index][lang][backend.GetName()].Translation = resp.GetText()
						translations[index][lang][backend.GetName()].Time = time.Since(t) / time.Millisecond
						mu.Unlock()
					}

				}(s, l, i, b)
			}
		}
	}
	wg.Wait()
	return translations
}

func (t *TranslateAdapter) doRequestsTranslators(text string, languages []string, c chan<- *RawTranslationData) {

	wg := &sync.WaitGroup{}
	for _, v := range languages {
		for _, back := range t.translateBackend {
			wg.Add(1)
			go func(text, lang string, backend backend_full.IBackendFull) {
				defer wg.Done()
				t := time.Now()
				resp := backend.TranslateFull(text, lang)

				raw := &RawTranslationData{
					Source:      resp.GetSource(),
					Lang:        resp.GetLang(),
					Name:        backend.GetName(),
					Translation: resp.GetText(),
					Time:        time.Since(t) / time.Millisecond,
				}
				c <- raw
			}(text, v, back)
		}

	}
	wg.Wait()
	close(c)
}

func (t *TranslateAdapter) createContainer(text string, langs []string) *TranslationContainer {
	container := &TranslationContainer{
		Translations: TranslationBag{},
		Langs:        langs,
		Original:     text,
		RawTransData: make(map[string]map[string]*RawTranslationData),
	}

	for _, lang := range container.Langs {
		container.RawTransData[lang] = make(map[string]*RawTranslationData)
		for _, back := range t.translateBackend {
			container.RawTransData[lang][back.GetName()] = &RawTranslationData{Name: back.GetName()}
		}
		container.RawTransData[lang]["uni"] = &RawTranslationData{Name: "uni"}
	}
	return container
}

func array_uniq(langs []string) []string {
	set := make(map[string]struct{})
	for _, val := range langs {
		if val == "" {
			continue
		}
		set[val] = struct{}{}
	}
	var languages []string
	for lang, _ := range set {
		languages = append(languages, lang)
	}
	return languages
}
