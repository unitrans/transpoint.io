// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package main
import (
	"net/http"
	"fmt"
	"regexp"
	"net/url"
	"bytes"
	_"io/ioutil"
)

func main(){
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://zafontan.ru", nil)
	res, err := client.Do(r)
	sessCookie := res.Header.Get("Set-Cookie")
	patt := regexp.MustCompile(`[\w\=]+`)
    sessionCookie := patt.FindString(sessCookie)

	form := &url.Values{}
	form.Add("question", "1")
	form.Add("module", "votes")
	form.Add("form_tag", "votes1")
	form.Add("result", "")
	form.Add("answer", "2")
	form.Add("ajax", "1")
	f := form.Encode()

	brut, _ := http.NewRequest("POST", "http://zafontan.ru", bytes.NewBufferString(f))


	brut.Header.Add("Cookie", sessionCookie + "; _ym_isad=1")
	fmt.Println(brut)
	res, err = client.Do(brut)
	fmt.Println(res, err)
//	body, _ := ioutil.ReadAll(res.Body)
//	fmt.Println(string(body))

}

