// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package main
import (
	"net/http"
	"fmt"
	"net/url"
	"bytes"
	_"io/ioutil"
	"time"
"math/big"

)

var client *http.Client

func main(){
	client = &http.Client{}

	for i := 0; i < 200; i++{
		time.Sleep(500 * time.Millisecond)
		brut()
		fmt.Println("Done", i)
	}




}

func brut() {
	r, _ := http.NewRequest("GET", "http://zafontan.ru", nil)
	res, err := client.Do(r)
	sessCookie := res.Header.Get("Set-Cookie")

	form := &url.Values{}
	form.Add("question", "1")
	form.Add("module", "votes")
	form.Add("form_tag", "votes1")
	form.Add("result", "")
	form.Add("answer", "2")
	form.Add("ajax", "1")
	f := form.Encode()

	brut, _ := http.NewRequest("POST", "http://zafontan.ru", bytes.NewBufferString(f))


	brut.Header.Add("Cookie", sessCookie)
	brut.Header.Add("Origin", "http://zafontan.ru")
	brut.Header.Add("Host", "http://zafontan.ru")
	brut.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	brut.Header.Add("X-Requested-With", "XMLHttpRequest")
	fmt.Println(brut)
	res, err = client.Do(brut)
	fmt.Println(res, err)
}

func factorial(n *big.Int) (result *big.Int) {
  //fmt.Println("n = ", n)
  b := big.NewInt(0)
  c := big.NewInt(1)

  if n.Cmp(b) == -1 {
    result = big.NewInt(1)
  }
  if n.Cmp(b) == 0 {
    result = big.NewInt(1)
  } else {
    // return n * factorial(n - 1);
    //fmt.Println("n = ", n)
    result = new(big.Int)
    result.Set(n)
    result.Mul(result, factorial(n.Sub(n, c)))
  }
  return
}

var X = `{
  "log": {
    "version": "1.2",
    "creator": {
      "name": "WebInspector",
      "version": "537.36"
    },
    "pages": [
      {
        "startedDateTime": "2015-12-16T14:38:57.809Z",
        "id": "page_18",
        "title": "http://zafontan.ru/",
        "pageTimings": {
          "onContentLoad": 2400.330066680908,
          "onLoad": 4457.499980926514
        }
      }
    ],
    "entries": [
      {
        "startedDateTime": "2015-12-16T14:39:19.576Z",
        "time": 114.65001106262207,
        "request": {
          "method": "POST",
          "url": "http://zafontan.ru/",
          "httpVersion": "HTTP/1.1",
          "headers": [
            {
              "name": "Cookie",
              "value": "_ym_isad=1; _ym_uid=1450220437752094722; SESSe0d7a4e680404d6f3072f9e0b0f5d94f=gs7hahidrvgirn23me39av29d2"
            },
            {
              "name": "Origin",
              "value": "http://zafontan.ru"
            },
            {
              "name": "Accept-Encoding",
              "value": "gzip, deflate"
            },
            {
              "name": "Host",
              "value": "zafontan.ru"
            },
            {
              "name": "Accept-Language",
              "value": "ru,en;q=0.8"
            },
            {
              "name": "User-Agent",
              "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 YaBrowser/15.6.2311.3451 (beta) Safari/537.36"
            },
            {
              "name": "Content-Type",
              "value": "application/x-www-form-urlencoded; charset=UTF-8"
            },
            {
              "name": "Accept",
              "value": "*/*"
            },
            {
              "name": "Referer",
              "value": "http://zafontan.ru/"
            },
            {
              "name": "X-Requested-With",
              "value": "XMLHttpRequest"
            },
            {
              "name": "Connection",
              "value": "keep-alive"
            },
            {
              "name": "Content-Length",
              "value": "63"
            }
          ],
          "queryString": [],
          "cookies": [
            {
              "name": "_ym_isad",
              "value": "1",
              "expires": null,
              "httpOnly": false,
              "secure": false
            },
            {
              "name": "_ym_uid",
              "value": "1450220437752094722",
              "expires": null,
              "httpOnly": false,
              "secure": false
            },
            {
              "name": "SESSe0d7a4e680404d6f3072f9e0b0f5d94f",
              "value": "gs7hahidrvgirn23me39av29d2",
              "expires": null,
              "httpOnly": false,
              "secure": false
            }
          ],
          "headersSize": 594,
          "bodySize": 63,
          "postData": {
            "mimeType": "application/x-www-form-urlencoded; charset=UTF-8",
            "text": "question=1&module=votes&form_tag=votes1&result=&answer=2&ajax=1",
            "params": [
              {
                "name": "question",
                "value": "1"
              },
              {
                "name": "module",
                "value": "votes"
              },
              {
                "name": "form_tag",
                "value": "votes1"
              },
              {
                "name": "result",
                "value": ""
              },
              {
                "name": "answer",
                "value": "2"
              },
              {
                "name": "ajax",
                "value": "1"
              }
            ]
          }
        },
        "response": {
          "status": 200,
          "statusText": "OK",
          "httpVersion": "HTTP/1.1",
          "headers": [
            {
              "name": "Date",
              "value": "Wed, 16 Dec 2015 14:39:19 GMT"
            },
            {
              "name": "Content-Encoding",
              "value": "gzip"
            },
            {
              "name": "Last-Modified",
              "value": "Thu, 03 Sep 2015 10:22:04 GMT"
            },
            {
              "name": "Server",
              "value": "Apache/2.2.23 (CentOS)"
            },
            {
              "name": "X-Powered-By",
              "value": "PHP/5.2.17"
            },
            {
              "name": "Vary",
              "value": "Accept-Encoding"
            },
            {
              "name": "Content-Type",
              "value": "text/html; charset=utf-8"
            },
            {
              "name": "Cache-Control",
              "value": "private, max-age=10800, pre-check=10800"
            },
            {
              "name": "Connection",
              "value": "keep-alive"
            },
            {
              "name": "Content-Length",
              "value": "1375"
            }
          ],
          "cookies": [],
          "content": {
            "size": 6766,
            "mimeType": "text/html",
            "compression": 5391
          },
          "redirectURL": "",
          "headersSize": 349,
          "bodySize": 1375,
          "_transferSize": 1724
        },
        "cache": {},
        "timings": {
          "blocked": 2.84100000135368,
          "dns": -1,
          "connect": -1,
          "send": 0.06800000119255989,
          "wait": 110.70099999778877,
          "receive": 1.0400110622870642,
          "ssl": -1
        },
        "connection": "112150",
        "pageref": "page_18"
      },
      {
        "startedDateTime": "2015-12-16T14:39:23.599Z",
        "time": 123.74997138977051,
        "request": {
          "method": "POST",
          "url": "http://zafontan.ru/",
          "httpVersion": "HTTP/1.1",
          "headers": [
            {
              "name": "Cookie",
              "value": "_ym_isad=1; _ym_uid=1450220437752094722; SESSe0d7a4e680404d6f3072f9e0b0f5d94f=gs7hahidrvgirn23me39av29d2"
            },
            {
              "name": "Origin",
              "value": "http://zafontan.ru"
            },
            {
              "name": "Accept-Encoding",
              "value": "gzip, deflate"
            },
            {
              "name": "Host",
              "value": "zafontan.ru"
            },
            {
              "name": "Accept-Language",
              "value": "ru,en;q=0.8"
            },
            {
              "name": "User-Agent",
              "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 YaBrowser/15.6.2311.3451 (beta) Safari/537.36"
            },
            {
              "name": "Content-Type",
              "value": "application/x-www-form-urlencoded; charset=UTF-8"
            },
            {
              "name": "Accept",
              "value": "*/*"
            },
            {
              "name": "Referer",
              "value": "http://zafontan.ru/"
            },
            {
              "name": "X-Requested-With",
              "value": "XMLHttpRequest"
            },
            {
              "name": "Connection",
              "value": "keep-alive"
            },
            {
              "name": "Content-Length",
              "value": "63"
            }
          ],
          "queryString": [],
          "cookies": [
            {
              "name": "_ym_isad",
              "value": "1",
              "expires": null,
              "httpOnly": false,
              "secure": false
            },
            {
              "name": "_ym_uid",
              "value": "1450220437752094722",
              "expires": null,
              "httpOnly": false,
              "secure": false
            },
            {
              "name": "SESSe0d7a4e680404d6f3072f9e0b0f5d94f",
              "value": "gs7hahidrvgirn23me39av29d2",
              "expires": null,
              "httpOnly": false,
              "secure": false
            }
          ],
          "headersSize": 594,
          "bodySize": 63,
          "postData": {
            "mimeType": "application/x-www-form-urlencoded; charset=UTF-8",
            "text": "question=1&module=votes&form_tag=votes1&result=&answer=2&ajax=1",
            "params": [
              {
                "name": "question",
                "value": "1"
              },
              {
                "name": "module",
                "value": "votes"
              },
              {
                "name": "form_tag",
                "value": "votes1"
              },
              {
                "name": "result",
                "value": ""
              },
              {
                "name": "answer",
                "value": "2"
              },
              {
                "name": "ajax",
                "value": "1"
              }
            ]
          }
        },
        "response": {
          "status": 200,
          "statusText": "OK",
          "httpVersion": "HTTP/1.1",
          "headers": [
            {
              "name": "Date",
              "value": "Wed, 16 Dec 2015 14:39:23 GMT"
            },
            {
              "name": "Content-Encoding",
              "value": "gzip"
            },
            {
              "name": "Last-Modified",
              "value": "Thu, 03 Sep 2015 10:22:04 GMT"
            },
            {
              "name": "Server",
              "value": "Apache/2.2.23 (CentOS)"
            },
            {
              "name": "X-Powered-By",
              "value": "PHP/5.2.17"
            },
            {
              "name": "Vary",
              "value": "Accept-Encoding"
            },
            {
              "name": "Content-Type",
              "value": "text/html; charset=utf-8"
            },
            {
              "name": "Cache-Control",
              "value": "private, max-age=10800, pre-check=10800"
            },
            {
              "name": "Connection",
              "value": "keep-alive"
            },
            {
              "name": "Content-Length",
              "value": "82"
            }
          ],
          "cookies": [],
          "content": {
            "size": 107,
            "mimeType": "text/html",
            "compression": 25
          },
          "redirectURL": "",
          "headersSize": 347,
          "bodySize": 82,
          "_transferSize": 429
        },
        "cache": {},
        "timings": {
          "blocked": 2.00100000074599,
          "dns": -1,
          "connect": -1,
          "send": 0.060000005760230124,
          "wait": 119.04399999912178,
          "receive": 2.644971384142508,
          "ssl": -1
        },
        "connection": "112150",
        "pageref": "page_18"
      },
      {
        "startedDateTime": "2015-12-16T14:39:33.628Z",
        "time": 27.72998809814453,
        "request": {
          "method": "POST",
          "url": "https://mc.yandex.ru/watch/26812653?page-url=http%3A%2F%2Fzafontan.ru%2F&browser-info=j%3A1%3As%3A1680x1050x24%3Ask%3A2%3Aadb%3A1%3Afpr%3A15455997001%3Acn%3A1%3Aw%3A1159x620%3Az%3A60%3Ai%3A20151216153933%3Aet%3A1450276774%3Aen%3Autf-8%3Av%3A672%3Ac%3A1%3Ala%3Aen-us%3Aar%3A1%3Anb%3A1%3Acl%3A242%3Als%3A214638960436%3Arqn%3A28%3Arn%3A488011773%3Ahid%3A629602460%3Ads%3A%2C%2C%2C%2C%2C%2C%2C%2C%2C4455%2C4455%2C4%2C%3Arqnl%3A1%3Ast%3A1450276774%3Au%3A1450220437752094722",
          "httpVersion": "HTTP/1.1",
          "headers": [
            {
              "name": "Origin",
              "value": "http://zafontan.ru"
            },
            {
              "name": "Accept-Encoding",
              "value": "gzip, deflate"
            },
            {
              "name": "Host",
              "value": "mc.yandex.ru"
            },
            {
              "name": "Accept-Language",
              "value": "ru,en;q=0.8"
            },
            {
              "name": "User-Agent",
              "value": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 YaBrowser/15.6.2311.3451 (beta) Safari/537.36"
            },
            {
              "name": "Content-Type",
              "value": "application/x-www-form-urlencoded"
            },
            {
              "name": "Accept",
              "value": "*/*"
            },
            {
              "name": "Referer",
              "value": "http://zafontan.ru/"
            },
            {
              "name": "Cookie",
              "value": "L=WFoEQwp5YVJ2YUIGf38ObQtsYkFmcHBHGD5BJD0uOAE3PiZA.1450111729.12101.333964.1848165e56680335e6c0cfdd231c9ce4; yandexuid=1129925841450220349; yabs-sid=769346901450220349; ys=wprid.1450219938975503-840207249112930312850268-iva1-0875"
            },
            {
              "name": "Connection",
              "value": "keep-alive"
            },
            {
              "name": "Content-Length",
              "value": "0"
            },
            {
              "name": "X-Yandex-Login",
              "value": "current=yk.multiship;last=yk.multiship"
            }
          ],
          "queryString": [
            {
              "name": "page-url",
              "value": "http%3A%2F%2Fzafontan.ru%2F"
            },
            {
              "name": "browser-info",
              "value": "j%3A1%3As%3A1680x1050x24%3Ask%3A2%3Aadb%3A1%3Afpr%3A15455997001%3Acn%3A1%3Aw%3A1159x620%3Az%3A60%3Ai%3A20151216153933%3Aet%3A1450276774%3Aen%3Autf-8%3Av%3A672%3Ac%3A1%3Ala%3Aen-us%3Aar%3A1%3Anb%3A1%3Acl%3A242%3Als%3A214638960436%3Arqn%3A28%3Arn%3A488011773%3Ahid%3A629602460%3Ads%3A%2C%2C%2C%2C%2C%2C%2C%2C%2C4455%2C4455%2C4%2C%3Arqnl%3A1%3Ast%3A1450276774%3Au%3A1450220437752094722"
            }
          ],
          "cookies": [
            {
              "name": "L",
              "value": "WFoEQwp5YVJ2YUIGf38ObQtsYkFmcHBHGD5BJD0uOAE3PiZA.1450111729.12101.333964.1848165e56680335e6c0cfdd231c9ce4",
              "expires": null,
              "httpOnly": false,
              "secure": false
            },
            {
              "name": "yandexuid",
              "value": "1129925841450220349",
              "expires": null,
              "httpOnly": false,
              "secure": false
            },
            {
              "name": "yabs-sid",
              "value": "769346901450220349",
              "expires": null,
              "httpOnly": false,
              "secure": false
            },
            {
              "name": "ys",
              "value": "wprid.1450219938975503-840207249112930312850268-iva1-0875",
              "expires": null,
              "httpOnly": false,
              "secure": false
            }
          ],
          "headersSize": 1173,
          "bodySize": 0
        },
        "response": {
          "status": 200,
          "statusText": "OK",
          "httpVersion": "HTTP/1.1",
          "headers": [
            {
              "name": "Pragma",
              "value": "no-cache"
            },
            {
              "name": "Date",
              "value": "Wed, 16 Dec 2015 14:39:33 GMT"
            },
            {
              "name": "Last-Modified",
              "value": "Wed, 16 Dec 2015 14:39:33 GMT"
            },
            {
              "name": "Server",
              "value": "nginx/1.8.0"
            },
            {
              "name": "Strict-Transport-Security",
              "value": "max-age=31536000"
            },
            {
              "name": "P3P",
              "value": "CP=\"NOI DEVa TAIa OUR BUS UNI STA\""
            },
            {
              "name": "Access-Control-Allow-Origin",
              "value": "http://zafontan.ru"
            },
            {
              "name": "Cache-Control",
              "value": "private, no-cache, no-store, must-revalidate, max-age=0"
            },
            {
              "name": "Access-Control-Allow-Credentials",
              "value": "true"
            },
            {
              "name": "Connection",
              "value": "keep-alive"
            },
            {
              "name": "Content-Type",
              "value": "image/gif"
            },
            {
              "name": "Content-Length",
              "value": "43"
            },
            {
              "name": "Expires",
              "value": "Wed, 16 Dec 2015 14:39:33 GMT"
            }
          ],
          "cookies": [],
          "content": {
            "size": 43,
            "mimeType": "image/gif",
            "compression": 0
          },
          "redirectURL": "",
          "headersSize": 497,
          "bodySize": 43,
          "_transferSize": 540
        },
        "cache": {},
        "timings": {
          "blocked": 3.50299999990966,
          "dns": -1,
          "connect": -1,
          "send": 0.0790000049164501,
          "wait": 22.68199999525679,
          "receive": 1.4659880980616329,
          "ssl": -1
        },
        "connection": "112743",
        "pageref": "page_18"
      }
    ]
  }
}`

