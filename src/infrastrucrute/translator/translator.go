// Copyright 2015 Yury Kozyrev. All rights reserved.
// Proprietary license.

// Package translator
package translator
import (
	"net/http"
	"time"
	"github.com/facebookgo/httpcontrol"
	"sync"
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