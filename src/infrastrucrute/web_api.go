// Copyright 2015 Home24 AG. All rights reserved.
// Proprietary license.
package infrastrucrute

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/urakozz/go-json-rest-middleware-jwt"
	"net/http"
	"time"
)

func NewWebApi() http.Handler {
	var jwt_middleware = &jwt.JWTMiddleware{
		Key:        []byte("testKey"),
		Realm:      "Unitrans",
		Timeout:    time.Second * 30,
		MaxRefresh: time.Hour * 24,
		Authenticator: func(userId string, password string) bool {
			return userId == "admin" && password == "admin"
		}}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "Authorization", "X-RequestId", "Origin"},
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


