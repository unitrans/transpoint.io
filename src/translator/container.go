// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translator

import (
	"github.com/urakozz/transpoint.io/src/translator/backend_particular"
	"time"
)

// TranslationBag hashmap
type TranslationBag map[string]string

// RawTransData struct
type RawTranslationData struct {
	Source      string
	Lang        string
	Name        string
	Translation string
	Time        time.Duration
	Score       float64
}

type RawParticularData struct {
	Original      string
	Source        string
	Lang          string
	Name          string
	ParticularBag []*ParticularItem
	Time          time.Duration
}

type ParticularItem struct {
	Url          string
	Order        int
	Original     string
	Time         time.Duration
	Translations []backend_particular.IParticularMeaning
}

// TranslationContainer struct
type TranslationContainer struct {
	Translations      TranslationBag
	Original          string
	Source            string
	//Meta              map[string]interface{}
	//RawTranslations   map[string]TranslationBag
	RawTransData      map[string]map[string]*RawTranslationData
	//RawParticularData []*RawParticularData
}
