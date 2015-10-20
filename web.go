// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package main

import (
	"net/http"
	"github.com/urakozz/transpoint.io/middleware"
	"html/template"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/goods/httpbuf"
	"github.com/satori/go.uuid"
	"github.com/OneOfOne/xxhash/native"


	"encoding/json"
	"strings"
	"log"
	"strconv"
	"fmt"
//	"regexp"
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
	Keys  []string    `json:"keys"`
}

func (u *User) IsLogin() bool {
	return u.Id != ""
}

type TemplateMap map[string]*template.Template
var Templates TemplateMap

func webInit() {
	parseFiles := func(name string) (*template.Template) {
		base := "templates/"
		partials := base + "partials/"
		mainTpl := base + name + ".html"
		return template.Must(template.ParseFiles(mainTpl, partials+"header.html", partials + "footer.html"))
	}
	Templates = make(TemplateMap)
	Templates["index"] = parseFiles("index")
	Templates["login"] = parseFiles("login")
	Templates["register"] = parseFiles("register")
	Templates["panel-index"] = parseFiles("panel-index")
}

func WebRouter() http.Handler {
	webInit()
	r := mux.NewRouter()

	r.Handle("/webapp", WebAction(WebIndex)).Methods("GET")
	r.Handle("/webapp/login", WebAction(WebLogin)).Methods("POST")
	r.Handle("/webapp/register", WebAction(WebRegisterGet)).Methods("GET")
	r.Handle("/webapp/register", WebAction(WebRegister)).Methods("POST")
	r.Handle("/webapp/logout", WebAction(WebLogout))
	r.Handle("/webapp/panel", WebAction(WebPanelIndex)).Methods("GET")
	r.Handle("/webapp/panel/keys", WebAction(WebPanelKeysPost)).Methods("POST")
	r.Handle("/webapp/panel/keys/delete/{id}", WebAction(WebPanelKeysDelete)).Methods("POST")

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
		if r.URL.Path == "/webapp" && context.Get(r, "user").(*User).IsLogin() {
			http.Redirect(w, r, "/webapp/panel", http.StatusFound)
			return
		}
		next(w, r)
	})
	app.UseHandler(r)
	return app
}

func WebIndex(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

//	var emojiiPattern = "[\\x{2712}\\x{2714}\\x{2716}\\x{271d}\\x{2721}\\x{2728}\\x{2733}\\x{2734}\\x{2744}\\x{2747}\\x{274c}\\x{274e}\\x{2753}-\\x{2755}\\x{2757}\\x{2763}\\x{2764}\\x{2795}-\\x{2797}\\x{27a1}\\x{27b0}\\x{27bf}\\x{2934}\\x{2935}\\x{2b05}-\\x{2b07}\\x{2b1b}\\x{2b1c}\\x{2b50}\\x{2b55}\\x{3030}\\x{303d}\\x{1f004}\\x{1f0cf}\\x{1f170}\\x{1f171}\\x{1f17e}\\x{1f17f}\\x{1f18e}\\x{1f191}-\\x{1f19a}\\x{1f201}\\x{1f202}\\x{1f21a}\\x{1f22f}\\x{1f232}-\\x{1f23a}\\x{1f250}\\x{1f251}\\x{1f300}-\\x{1f321}\\x{1f324}-\\x{1f393}\\x{1f396}\\x{1f397}\\x{1f399}-\\x{1f39b}\\x{1f39e}-\\x{1f3f0}\\x{1f3f3}-\\x{1f3f5}\\x{1f3f7}-\\x{1f4fd}\\x{1f4ff}-\\x{1f53d}\\x{1f549}-\\x{1f54e}\\x{1f550}-\\x{1f567}\\x{1f56f}\\x{1f570}\\x{1f573}-\\x{1f579}\\x{1f587}\\x{1f58a}-\\x{1f58d}\\x{1f590}\\x{1f595}\\x{1f596}\\x{1f5a5}\\x{1f5a8}\\x{1f5b1}\\x{1f5b2}\\x{1f5bc}\\x{1f5c2}-\\x{1f5c4}\\x{1f5d1}-\\x{1f5d3}\\x{1f5dc}-\\x{1f5de}\\x{1f5e1}\\x{1f5e3}\\x{1f5ef}\\x{1f5f3}\\x{1f5fa}-\\x{1f64f}\\x{1f680}-\\x{1f6c5}\\x{1f6cb}-\\x{1f6d0}\\x{1f6e0}-\\x{1f6e5}\\x{1f6e9}\\x{1f6eb}\\x{1f6ec}\\x{1f6f0}\\x{1f6f3}\\x{1f910}-\\x{1f918}\\x{1f980}-\\x{1f984}\\x{1f9c0}\\x{3297}\\x{3299}\\x{a9}\\x{ae}\\x{203c}\\x{2049}\\x{2122}\\x{2139}\\x{2194}-\\x{2199}\\x{21a9}\\x{21aa}\\x{231a}\\x{231b}\\x{2328}\\x{2388}\\x{23cf}\\x{23e9}-\\x{23f3}\\x{23f8}-\\x{23fa}\\x{24c2}\\x{25aa}\\x{25ab}\\x{25b6}\\x{25c0}\\x{25fb}-\\x{25fe}\\x{2600}-\\x{2604}\\x{260e}\\x{2611}\\x{2614}\\x{2615}\\x{2618}\\x{261d}\\x{2620}\\x{2622}\\x{2623}\\x{2626}\\x{262a}\\x{262e}\\x{262f}\\x{2638}-\\x{263a}\\x{2648}-\\x{2653}\\x{2660}\\x{2663}\\x{2665}\\x{2666}\\x{2668}\\x{267b}\\x{267f}\\x{2692}-\\x{2694}\\x{2696}\\x{2697}\\x{2699}\\x{269b}\\x{269c}\\x{26a0}\\x{26a1}\\x{26aa}\\x{26ab}\\x{26b0}\\x{26b1}\\x{26bd}\\x{26be}\\x{26c4}\\x{26c5}\\x{26c8}\\x{26ce}\\x{26cf}\\x{26d1}\\x{26d3}\\x{26d4}\\x{26e9}\\x{26ea}\\x{26f0}-\\x{26f5}\\x{26f7}-\\x{26fa}\\x{26fd}\\x{2702}\\x{2705}\\x{2708}-\\x{270d}\\x{270f}]|\\x{23}\\x{20e3}|\\x{2a}\\x{20e3}|\\x{30}\\x{20e3}|\\x{31}\\x{20e3}|\\x{32}\\x{20e3}|\\x{33}\\x{20e3}|\\x{34}\\x{20e3}|\\x{35}\\x{20e3}|\\x{36}\\x{20e3}|\\x{37}\\x{20e3}|\\x{38}\\x{20e3}|\\x{39}\\x{20e3}|\\x{1f1e6}[\\x{1f1e8}-\\x{1f1ec}\\x{1f1ee}\\x{1f1f1}\\x{1f1f2}\\x{1f1f4}\\x{1f1f6}-\\x{1f1fa}\\x{1f1fc}\\x{1f1fd}\\x{1f1ff}]|\\x{1f1e7}[\\x{1f1e6}\\x{1f1e7}\\x{1f1e9}-\\x{1f1ef}\\x{1f1f1}-\\x{1f1f4}\\x{1f1f6}-\\x{1f1f9}\\x{1f1fb}\\x{1f1fc}\\x{1f1fe}\\x{1f1ff}]|\\x{1f1e8}[\\x{1f1e6}\\x{1f1e8}\\x{1f1e9}\\x{1f1eb}-\\x{1f1ee}\\x{1f1f0}-\\x{1f1f5}\\x{1f1f7}\\x{1f1fa}-\\x{1f1ff}]|\\x{1f1e9}[\\x{1f1ea}\\x{1f1ec}\\x{1f1ef}\\x{1f1f0}\\x{1f1f2}\\x{1f1f4}\\x{1f1ff}]|\\x{1f1ea}[\\x{1f1e6}\\x{1f1e8}\\x{1f1ea}\\x{1f1ec}\\x{1f1ed}\\x{1f1f7}-\\x{1f1fa}]|\\x{1f1eb}[\\x{1f1ee}-\\x{1f1f0}\\x{1f1f2}\\x{1f1f4}\\x{1f1f7}]|\\x{1f1ec}[\\x{1f1e6}\\x{1f1e7}\\x{1f1e9}-\\x{1f1ee}\\x{1f1f1}-\\x{1f1f3}\\x{1f1f5}-\\x{1f1fa}\\x{1f1fc}\\x{1f1fe}]|\\x{1f1ed}[\\x{1f1f0}\\x{1f1f2}\\x{1f1f3}\\x{1f1f7}\\x{1f1f9}\\x{1f1fa}]|\\x{1f1ee}[\\x{1f1e8}-\\x{1f1ea}\\x{1f1f1}-\\x{1f1f4}\\x{1f1f6}-\\x{1f1f9}]|\\x{1f1ef}[\\x{1f1ea}\\x{1f1f2}\\x{1f1f4}\\x{1f1f5}]|\\x{1f1f0}[\\x{1f1ea}\\x{1f1ec}-\\x{1f1ee}\\x{1f1f2}\\x{1f1f3}\\x{1f1f5}\\x{1f1f7}\\x{1f1fc}\\x{1f1fe}\\x{1f1ff}]|\\x{1f1f1}[\\x{1f1e6}-\\x{1f1e8}\\x{1f1ee}\\x{1f1f0}\\x{1f1f7}-\\x{1f1fb}\\x{1f1fe}]|\\x{1f1f2}[\\x{1f1e6}\\x{1f1e8}-\\x{1f1ed}\\x{1f1f0}-\\x{1f1ff}]|\\x{1f1f3}[\\x{1f1e6}\\x{1f1e8}\\x{1f1ea}-\\x{1f1ec}\\x{1f1ee}\\x{1f1f1}\\x{1f1f4}\\x{1f1f5}\\x{1f1f7}\\x{1f1fa}\\x{1f1ff}]|\\x{1f1f4}\\x{1f1f2}|\\x{1f1f5}[\\x{1f1e6}\\x{1f1ea}-\\x{1f1ed}\\x{1f1f0}-\\x{1f1f3}\\x{1f1f7}-\\x{1f1f9}\\x{1f1fc}\\x{1f1fe}]|\\x{1f1f6}\\x{1f1e6}|\\x{1f1f7}[\\x{1f1ea}\\x{1f1f4}\\x{1f1f8}\\x{1f1fa}\\x{1f1fc}]|\\x{1f1f8}[\\x{1f1e6}-\\x{1f1ea}\\x{1f1ec}-\\x{1f1f4}\\x{1f1f7}-\\x{1f1f9}\\x{1f1fb}\\x{1f1fd}-\\x{1f1ff}]|\\x{1f1f9}[\\x{1f1e6}\\x{1f1e8}\\x{1f1e9}\\x{1f1eb}-\\x{1f1ed}\\x{1f1ef}-\\x{1f1f4}\\x{1f1f7}\\x{1f1f9}\\x{1f1fb}\\x{1f1fc}\\x{1f1ff}]|\\x{1f1fa}[\\x{1f1e6}\\x{1f1ec}\\x{1f1f2}\\x{1f1f8}\\x{1f1fe}\\x{1f1ff}]|\\x{1f1fb}[\\x{1f1e6}\\x{1f1e8}\\x{1f1ea}\\x{1f1ec}\\x{1f1ee}\\x{1f1f3}\\x{1f1fa}]|\\x{1f1fc}[\\x{1f1eb}\\x{1f1f8}]|\\x{1f1fd}\\x{1f1f0}|\\x{1f1fe}[\\x{1f1ea}\\x{1f1f9}]|\\x{1f1ff}[\\x{1f1e6}\\x{1f1f2}\\x{1f1fc}]";
//	rx := regexp.MustCompile(emojiiPattern)
//	var text = "a #💩 #and #🍦 #😳"
//	var i = -1
//	fmt.Println(rx.ReplaceAllStringFunc(text, func(s string) string {
//		i++
//		return strconv.Itoa(i)
//	}))
//	fmt.Fprint(w, text)
	Templates["login"].Execute(w, map[string]string{"Title":"Login", "token":ctx.CSRF})

	return
}

func WebIndexPage(w http.ResponseWriter, r *http.Request) {

	Templates["index"].Execute(w, nil)

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

		pHash := xxhash.Checksum64([]byte(pass))
		var passBytes []byte
		passBytes = strconv.AppendUint(passBytes, pHash, 10)
		if nil == err && user.Pass == string(passBytes) {
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
		pHash := xxhash.Checksum64([]byte(pass))
		var passBytes []byte
		passBytes = strconv.AppendUint(passBytes, pHash, 10)
		user := &User{Id:username, Pass:string(passBytes)}
		bytes, _ := json.Marshal(user)
		driver.Client.HSet("user", username, string(bytes))
		res, err = driver.Client.HGet("user", username).Result()
		ctx.Session.Values["user"] = username
		log.Println(res)
	}

	http.Redirect(w, r, "/webapp/panel", http.StatusFound)

	return
}

func WebRegisterGet(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

	Templates["register"].Execute(w, map[string]string{"Title":"Register", "token":ctx.CSRF})

	return
}

func WebLogout(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	delete(ctx.Session.Values, "user")

	http.Redirect(w, r, "/webapp", http.StatusFound)

	return
}

func WebPanelIndex(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

	keys, _ := driver.Client.HMGet("keys", ctx.User.Keys...).Result()
	keyMap := make(map[string]string, len(ctx.User.Keys))
	for i, key := range ctx.User.Keys {
		keyMap[key] = keys[i].(string)
	}

	Templates["panel-index"].Execute(w, map[string]interface{}{"Title":"Panel", "token":ctx.CSRF, "keys":keyMap, "user":ctx.User})

	return
}

func WebPanelKeysPost(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	uid := uuid.NewV4().String()
	uHash := xxhash.Checksum64([]byte(ctx.User.Id))
	var keyBytes []byte
	keyBytes = strconv.AppendUint(keyBytes, uHash, 10)
	keyBytes = append(keyBytes, '.')
	keyBytes = append(keyBytes, []byte(uid)[:8]...)

	var secretBytes []byte
	sHash := xxhash.Checksum64([]byte(uid))
	secretBytes = append(secretBytes, []byte(uid)[9:13]...)
	secretBytes = append(secretBytes, '.')
	secretBytes = strconv.AppendUint(secretBytes, sHash, 10)
	secretBytes = append(secretBytes, '.')
	secretBytes = append(secretBytes, []byte(uid)[24:]...)

	key := string(keyBytes)
	secret := string(secretBytes)
	ctx.User.Keys = append(ctx.User.Keys, key)
	bytes, _ := json.Marshal(ctx.User)
	driver.Client.HSet("user", ctx.User.Id, string(bytes))
	driver.Client.HSet("keys", key, secret)

	http.Redirect(w, r, "/webapp/panel", http.StatusFound)
	return
}

func WebPanelKeysDelete(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	key := mux.Vars(r)["id"]

	driver.Client.HDel("keys", key)

	for k, v := range ctx.User.Keys {
		if v == key {
			ctx.User.Keys = append(ctx.User.Keys[:k], ctx.User.Keys[k+1:]...)
		}
	}
	bytes, _ := json.Marshal(ctx.User)
	driver.Client.HSet("user", ctx.User.Id, string(bytes))

	http.Redirect(w, r, "/webapp/panel", http.StatusFound)
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