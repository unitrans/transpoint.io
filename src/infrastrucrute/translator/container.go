// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translator

import "time"

// TranslationBag hashmap
type TranslationBag map[string]string

// RawTransData struct
type RawTransData struct {
	Source      string
	Lang        string
	Name        string
	Translation string
	Time        time.Duration
}

// TranslationContainer struct
type TranslationContainer struct {
	Translations    TranslationBag
	Original        string
	Source          string
	Meta            map[string]interface{}
	RawTranslations map[string]TranslationBag
	RawTransData    []*RawTransData
}
