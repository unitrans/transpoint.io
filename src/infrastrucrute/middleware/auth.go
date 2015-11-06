// Copyright 2015 Yury Kozyrev. All rights reserved.
// Proprietary license.

// Package middleware - infrastructure/middleware
package middleware

import (
	"log"
	"net/http"

	e "github.com/pjebs/jsonerror"
	"gopkg.in/unrolled/render.v1"
)

// AuthConfig config for the middleware
type AuthConfig struct {
	ErrorMessages map[int]map[string]string
	Context       func(r *http.Request, authenticatedKey string)
	Debug         bool
	Client        func(key, secret string) bool

}

// RESTGate middleware
type RESTGate struct {
	headerKeyLabel    string
	headerSecretLabel string
	config            AuthConfig
}

// NewAuthMiddleware new Auth Middleware
func NewAuthMiddleware(headerKeyLabel string, headerSecretLabel string, config AuthConfig) *RESTGate {
	t := &RESTGate{headerKeyLabel: headerKeyLabel, headerSecretLabel: headerSecretLabel, config: config}

	if headerKeyLabel == "" { //headerKeyLabel must be defined
		if t.config.Debug == true {
			log.Printf("RestGate: headerKeyLabel is not defined.")
		}
		return nil
	}

	//Default Error Messages
	if t.config.ErrorMessages == nil {
		t.config.ErrorMessages = map[int]map[string]string{
			1:  e.New(1, "No Key Or Secret", "").Render(),
			2:  e.New(2, "Unauthorized Access", "").Render(),
		}
	} else {
		if _, ok := t.config.ErrorMessages[1]; !ok {
			t.config.ErrorMessages[1] = e.New(1, "No Key Or Secret", "").Render()
		}

		if _, ok := t.config.ErrorMessages[2]; !ok {
			t.config.ErrorMessages[2] = e.New(2, "Unauthorized Access", "").Render()
		}
	}

	return t
}

func (rest *RESTGate) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	//Check key in Header
	key := req.Header.Get(rest.headerKeyLabel)
	secret := req.Header.Get(rest.headerSecretLabel)

	if key == "" {
		//Authentication Information not included in request
		r := render.New(render.Options{})
		r.JSON(w, http.StatusUnauthorized, rest.config.ErrorMessages[1]) //"No Key Or Secret"
		return
	}

	authenticationPassed := rest.config.Client(key, secret)
	if authenticationPassed == false {
		r := render.New(render.Options{})
		r.JSON(w, http.StatusUnauthorized, rest.config.ErrorMessages[2]) //"Unauthorized Access"
		return
	}

	if rest.config.Context != nil {
		rest.config.Context(req, key)
	}
	next(w, req)

}
