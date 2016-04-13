// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package backend_particular

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const LG_URL = "http://www.lingvolive.com/api/Translation/Translate/"
const LG_PART_URL = "http://www.lingvolive.com/api/Translation/WordListPart/" //Fucking slow

var langsMap = map[string]string{
	"ru": "1049",
	"en": "1033",
	"da": "1030",
	"zh": "1028", //Chineese
	"nl": "1043",
	"fi": "1035",
	"de": "32775",
	"fr": "1036",
	"el": "1049", //Greek
	"hu": "1038", //Hungarian
	"it": "1040",
	"la": "1042", //Latin
	"no": "1044", //Norwegian
	"pl": "1045", //Polish
	"pt": "2070", //Portugal
	"es": "1034", //Spanish
	"tt": "1092", //Tatar
	"tr": "1055", //Turkish
	"uk": "1058", //Ukrainian
}

type AbbyyLingvoLiveTranslator struct {
	client *http.Client
}

func NewAbbyyLingvoLiveTranslator(c *http.Client) IBackendParticular{
	return &AbbyyLingvoLiveTranslator{client:c}
}

func (t *AbbyyLingvoLiveTranslator) TranslateWord(text string, language, to string) IBackendParticularResponse {

	data := &LingvoLiveTranslatorResponseFull{}
	data.Lang = language
	_, ok1 := langsMap[language]
	_, ok2 := langsMap[to]
	if !ok1 || !ok2 {
		return data
	}
	if language != "ru" && to != "ru" {
		return data
	}
	reqUrl := LG_URL+"?"+t.getQueryStringFull(text, language, to)
	req, _ := http.NewRequest("GET", reqUrl, nil)
	data.Url = reqUrl
	resp, err := t.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	reader := ioutil.NopCloser(resp.Body)

	// log.Println(str)
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		log.Println("error decode", err)
	}

	return data
}

func (t *AbbyyLingvoLiveTranslator) GetName() string {
	return "lingvo_live"
}

func (t *AbbyyLingvoLiveTranslator) getQueryStringFull(text, from, to string) string {
	form := url.Values{}
	form.Add("dstLang", langsMap[from])
	form.Add("srcLang", langsMap[to])
	form.Add("text", text)
	return form.Encode()
}

func (t *AbbyyLingvoLiveTranslator) getQueryStringPart(text, from, to string) string {
	form := url.Values{}
	form.Add("dstLang", from)
	form.Add("srcLang", to)
	form.Add("prefix", text)
	form.Add("pageSize", "10")
	form.Add("startIndex", "0")
	return form.Encode()
}

type LingvoLiveTranslatorResponseFull struct {
	Lang                  string
	Url                   string
	GlossaryUnits         interface{} `json:"glossaryUnits"`
	LanguagesReversed     bool        `json:"languagesReversed"`
	SeeAlsoWordForms      []string    `json:"seeAlsoWordForms"`
	Suggests              interface{} `json:"suggests"`
	WordByWordTranslation interface{} `json:"wordByWordTranslation"`
	Articles              []struct {
		Heading    string `json:"heading"`
		Dictionary string `json:"dictionary"`
		BodyHtml   string `json:"bodyHtml"`
	} `json:"lingvoArticles"`
}

type LingvoLiveTranslatorResponsePart struct {
	Lang string

	Items []struct {
		Heading      string `json:"heading"`
		Dictationary string `json:"lingvoDictionaryName"`
		Translations string `json:"lingvoTranslations"`
	} `json:"items"`
}

func (t *LingvoLiveTranslatorResponseFull) GetUrl() string {
	return t.Url
}

func (t *LingvoLiveTranslatorResponseFull) GetMeanings() []IParticularMeaning {
	meanings := []IParticularMeaning{}
	for _, v := range t.Articles {
		meaning := &Meaning{Dictionary: v.Dictionary}
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(v.BodyHtml))

		table := doc.Find(".article .article-body .article-body-items")
		table.Find(".article-body-items").Each(func(i int, s *goquery.Selection) {

			if s.Find(".paragraph-marker-top-level").Text() == "" {
				if s.Find(".parts-of-speech").Text() != "" && len(s.Find(".article-text").Nodes) == 0 {
					return
				}
			}

			value := s.Find(".article-text-wrap .article-text").Text()
			if value == "" {
				// maybe comment
				value = s.Find(".article-text-wrap .comment").Text()
			}
			value = strings.TrimLeft(value, "<-s, ->")
			value = strings.TrimLeft(value, ";")
			value = strings.TrimSpace(value)
			if "" != value {
				meaning.All = append(meaning.All, value)
			}
			if len(meaning.All) > 0 && meaning.Text == "" {

				meaning.Text = meaning.All[0]
			}
		})
		meanings = append(meanings, meaning)
	}
	return meanings
}
