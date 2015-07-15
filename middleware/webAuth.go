// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package middleware
import (
	"time"
	"net/http"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"github.com/justinas/nosurf"
	"github.com/goods/httpbuf"
	"encoding/json"
	"log"
)

type Session struct {
	session sessions.Store
}

func NewSession(session sessions.Store) *Session {
	return &Session{session}
}

func (a *Session) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := a.session.Get(r, "_session")
	if err != nil {
		http.SetCookie(rw, &http.Cookie{Name:"_session", Expires:time.Now().Add(-1*time.Hour)})
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print("zero ", session.Values)
	session.Values["time"] = time.Now().String()

	context.Set(r, "session", session)

	buf := new(httpbuf.Buffer)
	next(buf, r)

	session.Save(r, rw)
	buf.Apply(rw)

	//	res := rw.(ResponseWriter)
}

type CsrfMiddleware struct {
	name string
}

func NewCsrfMiddleware(name string) *CsrfMiddleware {
	return &CsrfMiddleware{name}
}

func (m *CsrfMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var token string
	var passed bool

	// nosurf disposes of the token as soon as it calls the http.Handler you provide...
	// in order to use it as negroni middleware, pull out token and dispose of it ourselves
	csrfHandler := nosurf.NewPure(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		token = nosurf.Token(r)
		passed = true
	}))

	csrfHandler.ServeHTTP(w, r)

	// csrf passed
	if passed {
		context.Set(r, m.name, token)
		next(w, r)
		context.Delete(r, m.name)
	} else {
		http.Error(w, "Invalid CSRF token", http.StatusBadRequest)
	}
}


type UserMiddleware struct {
	c *UserMiddlewareConfig
}

type UserMiddlewareConfig struct {
	Prototype interface{}
	Accessor  func(userId string) (string, error)
}

func NewUserMiddleware(c *UserMiddlewareConfig) *UserMiddleware {
	return &UserMiddleware{c}
}

func (m *UserMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := context.Get(r, "session").(*sessions.Session)
	userObj := m.c.Prototype
	if id, ok := s.Values["user"]; ok {
		if user, err := m.c.Accessor(id.(string)); err == nil {
			json.Unmarshal([]byte(user), &userObj)
		} else {
			http.Error(w, "User not found or data storage is not available", http.StatusBadRequest)
			return
		}
	}
	context.Set(r, "user", userObj)
	next(w, r)
}