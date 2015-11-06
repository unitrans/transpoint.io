// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package scenario

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/urakozz/go-json-rest-middleware-jwt"
	"net/http"
	"time"
	"log"
)

func NewWebApi() http.Handler {
	var jwt_middleware = &jwt.JWTMiddleware{
		Key:        []byte("testKey"),
		Realm:      "Unitrans",
		Timeout:    time.Second * 30,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			log.Println(userId, password)
//			return userId == "admin" && password == "admin"
			return password == "123123"
		}}

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
	)
	api.SetApp(api_router)
	return api.MakeHandler()
}



func handle_auth(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"authed": r.Env["REMOTE_USER"].(string)})
}


