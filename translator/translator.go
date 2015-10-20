// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package translator
import (
	"net/http"
	"time"
	"github.com/facebookgo/httpcontrol"
)

type Translator interface {
	Translate(text string, languages []string) TranslationBag
}


type TranslationBag map[string]string

type TranslationContainer struct {
	Translations TranslationBag
	Source       string
}

var client *http.Client

func initClient() (*http.Client) {
	if nil == client {
		transport := &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries: 3,
		}

		client = &http.Client{
			Transport: transport,
		}
	}
	return client
}