// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package middleware
import (
	"net/http"
	"github.com/gorilla/sessions"
	"log"
	"time"
	"github.com/gorilla/context"
)

type WebAuth struct {
	store sessions.Store
}
func NewWebAuth(store sessions.Store) *WebAuth {
	return &WebAuth{store}
}

func (a *WebAuth) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, err := a.store.Get(r, "_session")
	if err != nil {
		rw.WriteHeader(400)
		return
	}
	log.Print("zero ", session.Values)
	session.Values["time"] = time.Now().String()
	log.Print("before ", session.Values)
	context.Set(r, "session", session)
	next(rw, r)
	log.Print("after ", session.Values)
	session.Save(r, rw)

//	res := rw.(ResponseWriter)
}