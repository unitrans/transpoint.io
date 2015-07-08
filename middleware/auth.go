// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package middleware

/*
|--------------------------------------------------------------------------
| WARNING
|--------------------------------------------------------------------------
| Never Set HTTPSProtectionOff=true In Production.
| The Key and Password will be exposed and highly unsecure otherwise!
| The database server should also use HTTPS Connection and be hidden away
|
*/

/*
Thanks to Ido Ben-Natan ("IdoBn") for postgres fix.
Thanks to Jeremy Saenz & Brendon Murphy for timing-attack protection
*/

import (
	"log"
	"net/http"

	e "github.com/pjebs/jsonerror"
	"gopkg.in/unrolled/render.v1"
)

type AuthConfig struct {
	ErrorMessages map[int]map[string]string
	Context       func(r *http.Request, authenticatedKey string)
	Debug         bool
	Client        func(key, secret string) bool

}

type RESTGate struct {
	headerKeyLabel    string
	headerSecretLabel string
	config            AuthConfig
}

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

func (self *RESTGate) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	//Check key in Header
	key := req.Header.Get(self.headerKeyLabel)
	secret := req.Header.Get(self.headerSecretLabel)

	if key == "" {
		//Authentication Information not included in request
		r := render.New(render.Options{})
		r.JSON(w, http.StatusUnauthorized, self.config.ErrorMessages[1]) //"No Key Or Secret"
		return
	}

	authenticationPassed := self.config.Client(key, secret)
	if authenticationPassed == false {
		r := render.New(render.Options{})
		r.JSON(w, http.StatusUnauthorized, self.config.ErrorMessages[2]) //"Unauthorized Access"
		return
	}

	if self.config.Context != nil {
		self.config.Context(req, key)
	}
	next(w, req)

}
