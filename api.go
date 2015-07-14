// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.
package main

import (
	"fmt"
	"log"
	"time"

	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"github.com/satori/go.uuid"
	"github.com/codegangsta/negroni"

	"github.com/urakozz/transpoint.io/storage"
	"github.com/urakozz/transpoint.io/middleware"
)


const ApiVersion = "v1"

func ApiRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1", wrap(Default)).Methods("GET")
	r.HandleFunc("/v1/translations", wrap(Create)).Methods("POST")
	r.HandleFunc("/v1/translations/{id:[a-z0-9-]+}", wrap(Save)).Methods("POST")
	r.HandleFunc("/v1/translations/{id}", wrap(Get)).Methods("GET")
	r.HandleFunc("/v1/translations/{id}/{lang:[a-z]{2}}", wrap(GetParticular)).Methods("GET")
	r.HandleFunc("/v1/translations/{id}", wrap(Delete)).Methods("DELETE")
	r.HandleFunc("/v1/translations/{id}/{lang:[a-z]{2}}", wrap(DeleteParticular)).Methods("DELETE")

	app := negroni.New()
//	app.Use(negroni.NewRecovery())
//	app.Use(negroni.NewLogger())
	app.Use(middleware.NewAuthMiddleware(
		"X-Auth-Key",
		"X-Auth-Secret",
		middleware.AuthConfig{
			Context: func(r *http.Request, authenticatedKey string) {
				context.Set(r, 0, authenticatedKey)
			},
			Client: func(key, secret string) bool {
				sec := driver.Client.HGet("keys", key).Val()
				return sec == secret
			},
		},
	))
	app.UseHandler(r)

	return app
}

func ApiPing() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("PONG"))
	}
}

func wrap(action Action) (func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		v, code := action(w, r)

		callback := r.URL.Query().Get("callback")
		if callback == "" {
			render.JSON(w, code, v)
		} else {
			render.JSONP(w, code, callback, v)
		}
	}
}


func Default(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	methodMap := make(map[string]string)
	methodMap["translation_map"] = fmt.Sprintf("/%s/%s", ApiVersion, "translations")
	return methodMap, http.StatusOK
}

func Create(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	u1 := uuid.NewV4().String()
	id := context.Get(r, 0).(string) + "%" + u1
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	bag, _ := SmartSave(request, id)

	w.Header().Set("Location", "/"+ApiVersion+"/translations/"+u1)
	return bag, http.StatusCreated
}

func Save(w http.ResponseWriter, r *http.Request) (bag interface{}, status int) {
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	id := getId(r)
	bag, newLng := SmartSave(request, id)
	status = http.StatusOK

	if newLng > 0 {
		w.Header().Set("Location", "/"+ApiVersion+"/translations/"+id)
		status = http.StatusCreated
	}
	return
}

func SmartSave(request *RequestObject, id string) (bag storage.TranslationBag, newLng int) {
	bag, err := driver.GetAll(id)
	log.Println(bag, err)

	langs := request.Lang
	if err == nil {
		var newLangs []string
		if bag.Original == request.Text {
			for _, lang := range langs {
				if _, exists := bag.Translations[lang]; !exists {
					newLangs = append(newLangs, lang)
				}
			}
		} else {
			newLangs = langs
			for lang, _ := range bag.Translations {
				newLangs = append(newLangs, lang)
			}
		}

		langs = newLangs
	}

	newLng = len(langs)
	if 0 == newLng {
		return
	}
	container := translator.Translate(request.Text, langs)
	driver.Save(id, container.Source, request.Text, container.Translations)
	bag, err = driver.GetAll(id)
	log.Println(bag, err)
	return
}

func Get(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := getId(r)
	start := time.Now()
	bag, err := driver.GetAll(id)
	log.Printf("Completed in %v", time.Since(start))
	if nil != err {
		return map[string]string{"error":err.Error()}, http.StatusNotFound
	}

	return bag, http.StatusOK
}

func GetParticular(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := getId(r)
	lang := mux.Vars(r)["lang"]

	start := time.Now()
	bag, err := driver.GetLang(id, lang)
	log.Printf("Completed in %v", time.Since(start))

	if nil != err {
		return map[string]string{"error":err.Error()}, http.StatusNotFound
	}

	return bag, http.StatusOK
}

func Delete(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := getId(r)
	err := driver.Delete(id)
	if err != nil {
		log.Println("Delete", id, err)
	}
	return nil, http.StatusNoContent
}

func DeleteParticular(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := getId(r)
	lang := mux.Vars(r)["lang"]
	err := driver.DeleteLang(id, lang)

	if err != nil {
		log.Println("Delete", id, err)
	}
	return nil, http.StatusNoContent
}

func getId(r *http.Request) string {
	id := mux.Vars(r)["id"]
	key := context.Get(r, 0).(string)
	return key+"%"+id
}

type RequestObject struct {
	Id     string `json:"id"`
	Text   string `json:"text"`
	Lang   []string `json:"lang"`
	Source string `json:"source"`
}

