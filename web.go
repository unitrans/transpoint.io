// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package main

import (
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/urakozz/transpoint.io/middleware"
	"html/template"
	"time"
	"io"
	"crypto/rand"
	"encoding/base64"
	"crypto/subtle"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"github.com/goods/httpbuf"
	"log"
	"errors"
	"encoding/json"
)

const (
	TokenLength int = 32
	TtlDuration time.Duration = 20 * time.Minute
)

type WebContext struct {
	Session *sessions.Session
	User    *User
}
type WebAction func(w http.ResponseWriter, r *http.Request, ctx *WebContext) (error)

func (a WebAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := new(httpbuf.Buffer)
	ctx := &WebContext{
		Session: context.Get(r, "session").(*sessions.Session),
		User: context.Get(r, "user").(*User),
	}
	err := a(buf, r, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	buf.Apply(w)
}

type User struct {
	Id    string      `json:"id"`
	Email string      `json:"email"`
	Token string      `json:"token"`
	Ttl   time.Time   `json:"ttl"`
}

// RefreshToken refreshes Ttl and Token for the User.
func (u *User) RefreshToken() error {
	token := make([]byte, TokenLength)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return err
	}
	u.Token = base64.URLEncoding.EncodeToString(token)
	u.Ttl = time.Now().UTC().Add(TtlDuration)
	return nil
}

// IsValidToken returns a bool indicating that the User's current token hasn't
// expired and that the provided token is valid.
func (u *User) IsValidToken(token string) bool {
	if u.Ttl.Before(time.Now().UTC()) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(u.Token), []byte(token)) == 1
}
func (u *User) IsLogin() bool {
	return u.Id != ""
}

func WebRouter() http.Handler {
	r := mux.NewRouter()

	r.Handle("/webapp", WebAction(WebIndex)).Methods("GET")

	app := negroni.New()
	app.Use(middleware.NewSession(cookieStore))
	app.Use(middleware.NewCsrfMiddleware("csrf"))
	app.Use(middleware.NewUserMiddleware(&middleware.UserMiddlewareConfig{
		Accessor:func(userId string) (res string, err error) {
			if userId == "11q" {
				json, _ := json.Marshal(&User{Id:"11q", Email:"urakozz@me.com"})
				res = string(json)
				return
			}
			res, err = driver.Client.HGet("user", userId).Result()
			if res == "" {
				err = errors.New("not found")
			}
			return
		},
		Prototype: &User{},
	}))
	app.UseHandler(r)
	return app
}

func WebIndex(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	session := context.Get(r, "session").(*sessions.Session)
	csrf := context.Get(r, "csrf").(string)
	user := context.Get(r, "user")
	log.Println(session.Values)
	log.Println(csrf)
	log.Println(user)
	log.Println(ctx.User)
	log.Println(ctx.Session)
	session.Values["web"] = "some"
	session.Values["user"] = "11q"
	var homeTempl = template.Must(template.ParseFiles("templates/index.html"))
	homeTempl.Execute(w, r.Host)

	return
}




// UserMiddleware checks for the User in the session and adds them to the request context if they exist.
//func UserMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//	s := GetSession(r)
//	if id, ok := s.Values[sessionUser]; ok {
//		if user, err := dbmap.Get(User{}, id.(int64)); err == nil {
//			SetContextUser(user.(*User), r)
//		}
//	}
//	next(w, r)
//}
//
//// LoginRequiredMiddleware ensures a User is logged in, otherwise redirects them to the login page.
//func LoginRequiredMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//	if !IsLoggedIn(r) {
//		http.Redirect(w, r, "/", http.StatusFound)
//		return
//	}
//	next(w, r)
//}