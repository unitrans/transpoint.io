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
//	"io"
//	"crypto/rand"
//	"encoding/base64"
//	"crypto/subtle"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"github.com/goods/httpbuf"


	"encoding/json"
	"fmt"
	"strings"
	"log"
)

const (
	TokenLength int = 32
	TtlDuration time.Duration = 20 * time.Minute
)

type WebContext struct {
	Session *sessions.Session
	User    *User
	CSRF    string
}

type WebAction func(w http.ResponseWriter, r *http.Request, ctx *WebContext) (error)

func (a WebAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := new(httpbuf.Buffer)
	ctx := &WebContext{
		Session: context.Get(r, "session").(*sessions.Session),
		User: context.Get(r, "user").(*User),
		CSRF: context.Get(r, "csrf").(string),
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
	Pass  string      `json:"pass"`
}

//// RefreshToken refreshes Ttl and Token for the User.
//func (u *User) RefreshToken() error {
//	token := make([]byte, TokenLength)
//	if _, err := io.ReadFull(rand.Reader, token); err != nil {
//		return err
//	}
//	u.Token = base64.URLEncoding.EncodeToString(token)
//	u.Ttl = time.Now().UTC().Add(TtlDuration)
//	return nil
//}
//
//// IsValidToken returns a bool indicating that the User's current token hasn't
//// expired and that the provided token is valid.
//func (u *User) IsValidToken(token string) bool {
//	if u.Ttl.Before(time.Now().UTC()) {
//		return false
//	}
//	return subtle.ConstantTimeCompare([]byte(u.Token), []byte(token)) == 1
//}
func (u *User) IsLogin() bool {
	return u.Id != ""
}

type TemplateMap map[string]*template.Template
var Templates TemplateMap

func webInit(){
	parseFiles := func(name string)(*template.Template){
		base := "templates/"
		partials := "templates/partials/"
		mainTpl := base + name + ".html"
		return template.Must(template.ParseFiles(mainTpl, partials+"header.html",  partials + "footer.html"))
	}
	Templates = make(TemplateMap)
	Templates["login"] = parseFiles("login")
	Templates["panel-index"] = parseFiles("panel-index")
}

func WebRouter() http.Handler {
	webInit()
	r := mux.NewRouter()

	r.Handle("/webapp", WebAction(WebIndex)).Methods("GET")
	r.Handle("/webapp/login", WebAction(WebLogin)).Methods("POST")
	r.Handle("/webapp/register", WebAction(WebRegister)).Methods("POST")
	r.Handle("/webapp/logout", WebAction(WebLogout))
	r.Handle("/webapp/panel", WebAction(WebPanelIndex)).Methods("GET")

	app := negroni.New()
	app.Use(middleware.NewSession(cookieStore))
	app.Use(middleware.NewCsrfMiddleware("csrf"))
	app.Use(middleware.NewUserMiddleware(&middleware.UserMiddlewareConfig{
		Authenticator:func(userId string) (interface{}, error) {
			userObj := &User{}
			str, err := driver.Client.HGet("user", userId).Result()
			if err != nil {
				return nil, err
			}
			json.Unmarshal([]byte(str), &userObj)
			return userObj, nil
		},
	}))
	app.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if iface := context.Get(r, "user"); iface == nil {
			user := new(User)
			context.Set(r, "user", user)
		}
		next(w, r)
	})
	app.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if 0 == strings.Index(r.URL.Path, "/webapp/panel") && !context.Get(r, "user").(*User).IsLogin() {
			http.Redirect(w, r, "/webapp", http.StatusFound)
			return
		}
		if r.URL.Path == "/webapp" && context.Get(r, "user").(*User).IsLogin(){
			http.Redirect(w, r, "/webapp/panel", http.StatusFound)
			return
		}
		next(w, r)
	})
	app.UseHandler(r)
	return app
}

func WebIndex(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

	Templates["login"].Execute(w, map[string]string{"Title":"Some title", "Session": fmt.Sprintf("%+v", ctx.Session.Values), "token":ctx.CSRF, "user":fmt.Sprintf("%+v", ctx.User)})

	return
}

func WebLogin(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	session := ctx.Session

	r.ParseForm()
	username := r.Form.Get("username")
	pass := r.Form.Get("pass")

	res, err := driver.Client.HGet("user", username).Result()
	if res != "" && err == nil {
		user := &User{}
		err = json.Unmarshal([]byte(res), user)
		if nil == err && user.Pass == pass {
			session.Values["user"] = username
		}
	}

	http.Redirect(w, r, "/webapp/panel", http.StatusFound)

	return
}

func WebRegister(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

	r.ParseForm()
	username := r.Form.Get("username")
	pass := r.Form.Get("pass")

	res, err := driver.Client.HGet("user", username).Result()
	log.Println(res, err, "reg")
	if res == "" {
		user := &User{Id:username, Pass:pass}
		bytes, _ := json.Marshal(user)
		driver.Client.HSet("user", username, string(bytes))
		res, err = driver.Client.HGet("user", username).Result()
		log.Println(res)
	}

	http.Redirect(w, r, "/webapp/panel", http.StatusFound)

	return
}

func WebLogout(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	delete(ctx.Session.Values, "user")

	http.Redirect(w, r, "/webapp", http.StatusFound)

	return
}

func WebPanelIndex(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

	Templates["panel-index"].Execute(w, map[string]string{"Title":"Panel", "Session": fmt.Sprintf("%+v", ctx.Session.Values), "token":ctx.CSRF, "user":fmt.Sprintf("%+v", ctx.User)})

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