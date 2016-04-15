// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package backend_particular

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//import "gopkg.in/xmlpath.v2"

const MT_URL = "http://www.multitran.ru/c/m.exe"

var langsMap = map[string]string{"en": "1",
	"de":  "3",
	"fr":  "4",
	"es":  "5",
	"it":  "23",
	"nl":  "24",
	"lv":  "27",
	"et":  "26",
	"ja":  "28",
	"af":  "31",
	"eo":  "34",
	"xal": "35", //Kalmyk
}

type Multitran struct {
	client *http.Client
}

func NewMultitran(c *http.Client) IBackendParticular {
	return &Multitran{client: c}
}

func (t *Multitran) TranslateWord(text string, language, to string) IBackendParticularResponse {

	_, ok1 := langsMap[language]
	_, ok2 := langsMap[to]
	if !ok1 || !ok2 {
		return &MultitranResponse{}
	}

	reqUrl := LG_URL + "?" + t.getQueryStringFull(text, language, to)
	req, _ := http.NewRequest("GET", reqUrl, nil)
	resp, err := t.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
	}

	data := &MultitranResponse{}
	data.Lang = language
	data.Url = reqUrl
	data.BodyHtml = ioutil.ReadAll(resp.Body)

	return data
}

func (t *Multitran) GetName() string {
	return "lingvo_live"
}

func (t *Multitran) getQueryStringFull(text, from, to string) string {
	form := url.Values{}
	form.Add("l1", langsMap[from])
	form.Add("l2", langsMap[to])
	form.Add("s", text)
	//form.Add("CL", "1")
	return form.Encode()
}

type MultitranResponse struct {
	BodyHtml string
	Lang     string
	Url      string
}

func (t *MultitranResponse) GetUrl() string {
	return t.Url
}

func (t *MultitranResponse) GetMeanings() []IParticularMeaning {
	meanings := []IParticularMeaning{}
	//for _, v := range t.Articles {
		meaning := &Meaning{}
		//trs = html.xpath('//form[@id="translation"]/../table[2]/tr')
		//for tr in trs:
		//  tds = tr.xpath('td')
		//    for td in tds:
		//      for elem in td.xpath('descendant::text()'):
		//        translation += '%s' % elem.rstrip('\r\n')
		//        translation += '\n'
		//returnValue(translation)
		meanings = append(meanings, meaning)
	//}
	return meanings
}
