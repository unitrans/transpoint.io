// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package scenario

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/urakozz/go-json-rest-middleware-jwt"
	"net/http"
	"time"
	"encoding/json"
	"github.com/unitrans/unitrans/src/domain"
	"log"
)

type UserMiddleware struct {

}

func (mw *UserMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {

		userId := r.Env["REMOTE_USER"].(string)
		res, _ := userRepository.GetUserById(userId)
		user := domain.NewUser()
		json.Unmarshal([]byte(res), user)
		r.Env["USER"] = user

		// call the handler
		h(w, r)

	}
}

func NewWebApi() http.Handler {
	Authenticator := func(userId, password string) bool {

		res, err := userRepository.GetUserById(userId)
		if err != nil || res == "" {
			return false
		}
		user := domain.NewUser()
		err = json.Unmarshal([]byte(res), user)

		return nil == err && user.Pass == domain.HashPassword(password)
	}

	var jwt_middleware = &jwt.JWTMiddleware{
		Key:        []byte("testKey"),
		Realm:      "Unitrans",
		Timeout:    time.Hour * 30,
		MaxRefresh: time.Hour * 24,
		Authenticator: Authenticator,
	}

	var DevStack = []rest.Middleware{
		&rest.AccessLogApacheMiddleware{},
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.PoweredByMiddleware{
			XPoweredBy:"unitrans",
		},
		&rest.RecoverMiddleware{
			EnableResponseStackTrace: true,
		},
		&rest.JsonIndentMiddleware{},
	}
	api := rest.NewApi()

	api.Use(DevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login" && request.URL.Path != "/refresh" && request.URL.Path != "/register" && request.URL.Path != "/checkExists" && request.URL.Path != "/tr"
		},
		IfTrue: jwt_middleware,
	})
	api.Use(&rest.IfMiddleware{
		Condition: func(r *rest.Request) bool {
			_, ok := r.Env["REMOTE_USER"]
			return ok
		},
		IfTrue: &UserMiddleware{},
	})
	api_router, _ := rest.MakeRouter(
		rest.Post("/login", jwt_middleware.LoginHandler),
		rest.Get("/test", handle_auth),
		rest.Get("/refresh", jwt_middleware.RefreshHandler),
		rest.Post("/register", Register),
		rest.Post("/checkExists", CheckExists),

		rest.Get("/keys", KeysList),
		rest.Post("/keys", KeyCreate),
		rest.Delete("/keys/*key", KeyDelete),

		rest.Post("/tr", Translate),
	)
	api.SetApp(api_router)
	return api.MakeHandler()
}



func handle_auth(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"authed": r.Env["REMOTE_USER"].(string)})
}

func Register(w rest.ResponseWriter, r *rest.Request) {
	var form struct {
		Id string      `json:"username"`
		Pass  string   `json:"password"`
	}
	r.DecodeJsonPayload(&form)

	if form.Pass == "" {
		rest.Error(w, "Invalid input", 400)
		return
	}

	res, err := userRepository.GetUserById(form.Id)
	if res != "" {
		rest.Error(w, "[username]exists", 400)
		return
	}

	user := &domain.User{Id:form.Id, Pass:domain.HashPassword(form.Pass)}
	bytes, _ := json.Marshal(user)
	userRepository.SaveUserById(form.Id, string(bytes))
	res, err = userRepository.GetUserById(form.Id)
	log.Println(res, err)
	w.WriteJson(map[string]string{"Status":"OK"})
}

func CheckExists(w rest.ResponseWriter, r *rest.Request) {
	var form struct {
		Id string      `json:"username"`
		Pass  string   `json:"password"`
	}
	r.DecodeJsonPayload(&form)

	if form.Id == "" {
		rest.Error(w, "Invalid input", 400)
		return
	}

	res, _ := userRepository.GetUserById(form.Id)
	if res == "" {
		rest.Error(w, "[username] not exists", 400)
		return
	}

	w.WriteJson(map[string]string{"Status":"OK"})
}

func KeyCreate(w rest.ResponseWriter, r *rest.Request){

	user := r.Env["USER"].(*domain.User)

	key, secret := domain.GenerateKeyPair(user)
	user.AddKey(key)


	go func(user *domain.User, key, secret string){
		bytes, _ := json.Marshal(user)
		userRepository.SaveUserById(user.Id, string(bytes))
		userRepository.SaveSecretByKey(key, secret)
	}(user, key, secret)

	userClone := user.Clone()
	userClone.Pass = ""
	w.WriteHeader(http.StatusCreated)
	w.WriteJson(userClone)
}

func KeysList(w rest.ResponseWriter, r *rest.Request){

	user := r.Env["USER"].(*domain.User)

	user.Pass = ""
	w.WriteJson(user)
}

func KeyDelete(w rest.ResponseWriter, r *rest.Request){
	key := r.PathParam("key")
	user := r.Env["USER"].(*domain.User)

	for k, v := range user.Keys {
		if v == key {
			user.Keys = append(user.Keys[:k], user.Keys[k+1:]...)
			break
		}
	}
	if len(user.Keys) == 0 {
		user.Keys = nil
	}

	go func(user *domain.User, key string){
		bytes, _ := json.Marshal(user)
		userRepository.SaveUserById(user.Id, string(bytes))
		userRepository.DeleteSecretByKey(key)
	}(user, key)


	userClone := user.Clone()
	userClone.Pass = ""
	w.WriteJson(userClone)
}



func Translate(w rest.ResponseWriter, r *rest.Request){

	request := &RequestObject{}
	r.DecodeJsonPayload(request)
	c := translator.Translate(request.Text, request.Lang)
	w.WriteJson(c)
}

