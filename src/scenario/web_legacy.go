// Copyright 2015 Kozyrev Yury. All rights reserved.
// Proprietary license.
package scenario

import (
	"net/http"
	"html/template"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/goods/httpbuf"
	"github.com/OneOfOne/xxhash/native"


	"encoding/json"
	"strings"
	"log"
	"strconv"

	"github.com/unitrans/unitrans/src/infrastrucrute/middleware"
	"github.com/unitrans/unitrans/src/interface/repository/redis"
	"github.com/unitrans/unitrans/src/domain"
)

const (
	TokenLength int = 32
	TtlDuration time.Duration = 20 * time.Minute
)

var userRepository *repository.UserRepository

type WebContext struct {
	Session *sessions.Session
	User    *domain.User
	CSRF    string
}

type WebAction func(w http.ResponseWriter, r *http.Request, ctx *WebContext) (error)

func (a WebAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := new(httpbuf.Buffer)
	ctx := &WebContext{
		Session: context.Get(r, "session").(*sessions.Session),
		User: context.Get(r, "user").(*domain.User),
		CSRF: context.Get(r, "csrf").(string),
	}

	err := a(buf, r, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	buf.Apply(w)
}


var Templates = make(map[string]*template.Template)

func webInit() {
	parseFiles := func(name string) (*template.Template) {
		base := "templates/"
		partials := base + "partials/"
		mainTpl := base + name + ".html"
		return template.Must(template.ParseFiles(mainTpl, partials+"header.html", partials + "footer.html"))
	}

	Templates["index"] = parseFiles("index")
	Templates["login"] = parseFiles("login")
	Templates["register"] = parseFiles("register")
	Templates["panel-index"] = parseFiles("panel-index")
}

func WebRouter(u *repository.UserRepository, session sessions.Store) http.Handler {
	userRepository = u
	webInit()
	r := mux.NewRouter()

	r.Handle("/", WebAction(WebIndex)).Methods("GET")
	r.Handle("/login", WebAction(WebLogin)).Methods("POST")
	r.Handle("/register", WebAction(WebRegisterGet)).Methods("GET")
	r.Handle("/register", WebAction(WebRegister)).Methods("POST")
	r.Handle("/logout", WebAction(WebLogout))
	r.Handle("/panel", WebAction(WebPanelIndex)).Methods("GET")
	r.Handle("/panel/keys", WebAction(WebPanelKeysPost)).Methods("POST")
	r.Handle("/panel/keys/delete/{id}", WebAction(WebPanelKeysDelete)).Methods("POST")

	app := negroni.New()
	app.Use(middleware.NewSession(session))
	app.Use(middleware.NewCsrfMiddleware("csrf"))
	app.Use(middleware.NewUserMiddleware(&middleware.UserMiddlewareConfig{
		Authenticator:func(userId string) (interface{}, error) {
			userObj := domain.NewUser()
			str, err := userRepository.GetUserById(userId)
			if err != nil {
				return nil, err
			}
			json.Unmarshal([]byte(str), &userObj)
			return userObj, nil
		},
	}))
	app.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if iface := context.Get(r, "user"); iface == nil {
			user := domain.NewUser()
			context.Set(r, "user", user)
		}
		next(w, r)
	})
	app.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if 0 == strings.Index(r.URL.Path, "/panel") && !context.Get(r, "user").(*domain.User).IsLogin() {
			http.Redirect(w, r, "/webapp", http.StatusFound)
			return
		}
		if r.URL.Path == "/" && context.Get(r, "user").(*domain.User).IsLogin() {
			http.Redirect(w, r, "/webapp/panel", http.StatusFound)
			return
		}
		next(w, r)
	})
	app.UseHandler(r)
	return context.ClearHandler(app)
}

func WebIndex(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

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

	res, err := userRepository.GetUserById(username)
	if res != "" && err == nil {
		user := domain.NewUser()
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

	res, err := userRepository.GetUserById(username)
	log.Println(res, err, "reg")
	if res == "" && username != "" && pass != "" {
		pHash := xxhash.Checksum64([]byte(pass))
		var passBytes []byte
		passBytes = strconv.AppendUint(passBytes, pHash, 10)
		user := &domain.User{Id:username, Pass:string(passBytes)}
		bytes, _ := json.Marshal(user)
		userRepository.SaveUserById(username, string(bytes))
		res, err = userRepository.GetUserById(username)
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

	keys, _ := userRepository.GetAllSecretsByKeys(ctx.User.Keys...)
	keyMap := make(map[string]string, len(ctx.User.Keys))
	for i, key := range ctx.User.Keys {
		keyMap[key] = keys[i].(string)
	}
	log.Printf("%+v",ctx.User)

	Templates["panel-index"].Execute(w, map[string]interface{}{"Title":"Panel", "token":ctx.CSRF, "keys":keyMap, "user":ctx.User})

	return
}

func WebPanelKeysPost(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {

	key, secret := domain.GenerateKeyPair(ctx.User)

	ctx.User.AddKey(key)

	bytes, _ := json.Marshal(ctx.User)
	userRepository.SaveUserById(ctx.User.Id, string(bytes))
	userRepository.SaveSecretByKey(key, secret)

	http.Redirect(w, r, "/webapp/panel", http.StatusFound)
	return
}

func WebPanelKeysDelete(w http.ResponseWriter, r *http.Request, ctx *WebContext) (err error) {
	key := mux.Vars(r)["id"]

	userRepository.DeleteSecretByKey(key)

	for k, v := range ctx.User.Keys {
		if v == key {
			ctx.User.Keys = append(ctx.User.Keys[:k], ctx.User.Keys[k+1:]...)
		}
	}
	bytes, _ := json.Marshal(ctx.User)
	userRepository.SaveUserById(ctx.User.Id, string(bytes))

	http.Redirect(w, r, "/webapp/panel", http.StatusFound)
	return
}
