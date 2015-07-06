// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package translator
import (
	"net/http"
	"net"
	"time"
)

type Translator interface {
	Translate(text string, languages []string) TranslationBag
}


type TranslationBag map[string]string

type TranslationContainer struct {
	Translations TranslationBag
	Source string
}

var client *http.Client

func initClient() (*http.Client) {
	if nil == client {
		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: time.Minute,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}

		client = &http.Client{
			Transport: transport,
		}
	}
	return client
}