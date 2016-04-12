// Copyright 2016 Home24 AG. All rights reserved.
// Proprietary license.
package httpclient

import (
	"github.com/facebookgo/httpcontrol"
	"net/http"
	"time"
)

func GetHttpClient() *http.Client {
	transport := &httpcontrol.Transport{
		RequestTimeout:      time.Minute,
		MaxTries:            3,
		MaxIdleConnsPerHost: 10,
	}

	return &http.Client{
		Transport: transport,
	}
}
