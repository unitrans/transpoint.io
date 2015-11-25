// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package translator



// TranslationBag hashmap
type TranslationBag map[string]string

// TranslationContainer struct
type TranslationContainer struct {
	Translations TranslationBag
	Original string
	Source       string
	Meta         map[string]interface{}
	RawTranslations map[string]TranslationBag
}