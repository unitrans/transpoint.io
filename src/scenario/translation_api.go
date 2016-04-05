// Copyright 2015 Yury Kozyrev. All rights reserved.
// Proprietary license.
package scenario

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
	"github.com/rs/cors"


	r "gopkg.in/unrolled/render.v1"

	"github.com/urakozz/transpoint.io/src/infrastrucrute/middleware"
	"github.com/urakozz/transpoint.io/src/interface/repository/redis"
	t "github.com/urakozz/transpoint.io/src/infrastrucrute/translator"
)


const ApiVersion = "v1"

var (
	transRepository *repository.TranslationRepository
	translator t.Translator
	render *r.Render
)


type Action func(w http.ResponseWriter, r *http.Request) (interface{}, int)


func ApiRouter(userRepository *repository.UserRepository, tRep *repository.TranslationRepository, tr t.Translator) http.Handler {
	transRepository = tRep
	translator = tr

	render = r.New(r.Options{})

	r := mux.NewRouter()
	r.HandleFunc("/", wrap(Default)).Methods("GET")
	r.HandleFunc("/translations", wrap(ApiTranslate)).Methods("GET")
	r.HandleFunc("/translations", wrap(Create)).Methods("POST")
	r.HandleFunc("/translations/{id:[a-z0-9-]+}", wrap(Save)).Methods("POST")
	r.HandleFunc("/translations/{id}", wrap(Get)).Methods("GET")
	r.HandleFunc("/translations/{id}/{lang:[a-z]{2}}", wrap(GetParticular)).Methods("GET")
	r.HandleFunc("/translations/{id}", wrap(Delete)).Methods("DELETE")
	r.HandleFunc("/translations/{id}/{lang:[a-z]{2}}", wrap(DeleteParticular)).Methods("DELETE")

	app := negroni.New()
	//	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(middleware.NewAuthMiddleware(
		"X-Auth-Key",
		"X-Auth-Secret",
		middleware.AuthConfig{
			Context: func(r *http.Request, authenticatedKey string) {
				context.Set(r, 0, authenticatedKey)
			},
			Client: func(key, secret string) bool {
//				sec := userRepository.GetSecretByKey(key)
				log.Println(key, secret)
				return true
			},
		},
	))
	c := cors.New(cors.Options{})
	app.Use(c)
	app.UseHandler(r)

	return context.ClearHandler(app)
}

func ApiPing() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	}
}

func wrap(action Action) (func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		v, code := action(w, r)
//		log.Println(r)

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

func ApiTranslate(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	container := translator.Translate(r.URL.Query().Get("text"), r.URL.Query()["lang"])
	return container, http.StatusOK
}

func Create(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	u1 := uuid.NewV4().String()
	id := context.Get(r, 0).(string) + "%" + u1
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	bag, _ := SmartSave(request, id)

	w.Header().Set("Location", "/" + ApiVersion + "/translations/" + u1)
	return bag, http.StatusCreated
}

func Save(w http.ResponseWriter, r *http.Request) (bag interface{}, status int) {
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	id := getId(r)
	bag, newLng := SmartSave(request, id)
	status = http.StatusOK

	if newLng > 0 {
		w.Header().Set("Location", "/" + ApiVersion + "/translations/" + id)
		status = http.StatusCreated
	}
	return
}

func SmartSave(request *RequestObject, id string) (bag repository.TranslationBag, newLng int) {
	bag, err := transRepository.GetAll(id)
	log.Println(err)

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

	meta := make(map[string]interface{})
	container := translator.Translate(request.Text, langs)
	meta["raw"] = container.RawTranslations


	transRepository.Save(id, container.Source, request.Text, container.Translations, meta)
	bag, err = transRepository.GetAll(id)
//	log.Println(bag, err)
	return
}

func Get(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := getId(r)
	start := time.Now()
	bag, err := transRepository.GetAll(id)
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
	bag, err := transRepository.GetLang(id, lang)
	log.Printf("Completed in %v", time.Since(start))

	if nil != err {
		return map[string]string{"error":err.Error()}, http.StatusNotFound
	}

	return bag, http.StatusOK
}

func Delete(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := getId(r)
	err := transRepository.Delete(id)
	if err != nil {
		log.Println("Delete", id, err)
	}
	return nil, http.StatusNoContent
}

func DeleteParticular(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := getId(r)
	lang := mux.Vars(r)["lang"]
	err := transRepository.DeleteLang(id, lang)

	if err != nil {
		log.Println("Delete", id, err)
	}
	return nil, http.StatusNoContent
}

func getId(r *http.Request) string {
	id := mux.Vars(r)["id"]
	key := context.Get(r, 0).(string)
	return key + "%" + id
}

type RequestObject struct {
	Id     string `json:"id"`
	Text   string `json:"text"`
	Lang   []string `json:"lang"`
	Source string `json:"source"`
}

