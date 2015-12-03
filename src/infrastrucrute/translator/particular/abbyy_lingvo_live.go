// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package particular

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
//	"da": "1030",?
//	"zh": "1028",?
//	"nl": "1035",?
//	"fi": "1043",?
	"de": "32775",
	"fr": "1036",
}

type AbbyyLingvoLiveTranslator struct {
	client *http.Client
}

func (t *AbbyyLingvoLiveTranslator) SetClient(c *http.Client) {
	t.client = c
}

func (t *AbbyyLingvoLiveTranslator) TranslateOne(text string, language, to string) IParticularResponse {

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
	req, _ := http.NewRequest("GET", LG_URL+"?"+t.getQueryStringFull(text, language, to), nil)

	resp, err := t.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}
	reader := ioutil.NopCloser(resp.Body)
		b , _ := ioutil.ReadAll(reader)
	str := string(b)

		log.Println(str)
	if err := json.NewDecoder(strings.NewReader(str)).Decode(&data); err != nil {
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

func (t *LingvoLiveTranslatorResponseFull) GetMeanings() []IParticularMeaning {
	meanings := []IParticularMeaning{}
	for _, v := range t.Articles {
		meaning := &Meaning{Dictationary: v.Dictionary}
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(v.BodyHtml))

		table := doc.Find(".article .article-body .article-body-items")
		table.Find(".article-body-items .article-text-wrap .article-text").Each(func(i int, s *goquery.Selection) {
			value := strings.TrimSpace(s.Text())
			value = strings.Trim(value, ";")
			value = strings.Trim(value, ",")
			if len(value) > 0 && value[0] =='<' {
				value = ""
			}
			if "" != value {
				meaning.All = append(meaning.All, value)
			}
			if len(meaning.All) > 0 {
				meaning.Text = meaning.All[0]
			}
		})
		meanings = append(meanings, meaning)
	}
	return meanings
}
