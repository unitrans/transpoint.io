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
)

const (
	TokenLength int           = 32
	TtlDuration time.Duration = 20 * time.Minute
)

type WebContext struct {
	Session  *sessions.Session
	User     *User
}
type WebAction func(w http.ResponseWriter, r *http.Request, ctx *WebContext) (error)

func (a WebAction) ServeHTTP(w http.ResponseWriter, r *http.Request){
	buf := new(httpbuf.Buffer)
	ctx := &WebContext{
		Session: context.Get(r, "session").(*sessions.Session),
	}
	err := a(buf, r, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ctx.Session.Save(r, w)
	buf.Apply(w)
}

type User struct {
	Id        int64       `db:"id"`
	Email     string      `db:"email"`
	Token     string      `db:"token"`
	Ttl       time.Time   `db:"ttl"`
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

func WebRouter() http.Handler {
	r := mux.NewRouter()

	r.Handle("/webapp", WebAction(WebIndex)).Methods("GET")

	app := negroni.New()
	app.Use(middleware.NewWebAuth(cookieStore))
	app.UseHandler(r)
	return app
}

func WebIndex(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	session := context.Get(r, "session").(*sessions.Session)
	session.Values["web"] = "some"
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