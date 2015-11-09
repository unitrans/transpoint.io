// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package scenario

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/urakozz/go-json-rest-middleware-jwt"
	"net/http"
	"time"
	"encoding/json"
	"github.com/urakozz/transpoint.io/src/domain"
	"log"
)

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
		Timeout:    time.Second * 30,
		MaxRefresh: time.Hour * 24,
		Authenticator: Authenticator,
	}

	var DevStack = []rest.Middleware{
		&rest.AccessLogApacheMiddleware{},
		&rest.TimerMiddleware{},
		&rest.RecorderMiddleware{},
		&rest.PoweredByMiddleware{},
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
		AllowedMethods: []string{"*"},
		AllowedHeaders: []string{"*"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})
	api.Use(&rest.IfMiddleware{
		Condition: func(request *rest.Request) bool {
			return request.URL.Path != "/login" && request.URL.Path != "/refresh"
		},
		IfTrue: jwt_middleware,
	})
	api_router, _ := rest.MakeRouter(
		rest.Post("/login", jwt_middleware.LoginHandler),
		rest.Get("/test", handle_auth),
		rest.Get("/refresh", jwt_middleware.RefreshHandler),
		rest.Get("/register", Register),
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

	res, err := userRepository.GetUserById(form.Id)
	if err != nil || res == "" || form.Pass == "" {
		rest.Error(w, "Invalid input", 400)
	}

	user := &domain.User{Id:form.Id, Pass:domain.HashPassword(form.Pass)}
	bytes, _ := json.Marshal(user)
	userRepository.SaveUserById(form.Id, string(bytes))
	res, err = userRepository.GetUserById(form.Id)
	log.Println(res, err)
	w.WriteJson(map[string]string{"Status":"OK"})
}


