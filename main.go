package main

import (
	"os"
	"log"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/joho/godotenv"
//	"github.com/xlab/closer"
//	"github.com/pjebs/restgate"
	r "gopkg.in/unrolled/render.v1"

	"./storage"
	"github.com/satori/go.uuid"
	"encoding/json"
	"time"
	t "./translator"
)

const ApiVersion = "v1"

var (
	render *r.Render
	router *mux.Router
	driver *storage.RedisDriver
	translator *t.YandexTranslator
)

type Action func(w http.ResponseWriter, r *http.Request) (interface{}, int)

func init() {
	godotenv.Load()
	render = r.New(r.Options{})
	router = mux.NewRouter().StrictSlash(true)
	driver = storage.NewRedisDriver(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_DB"), os.Getenv("REDIS_PASS"))
	translator = t.NewYandexTranslator()
}

func main() {
	//	closer.Bind(cleanup)
	//	closer.Checked(run, true)
	run()
}

func run() error {
	port := os.Getenv("PORT")

	router.HandleFunc("/ping", wrap(Ping)).Methods("GET")
	router.HandleFunc("/"+ApiVersion, wrap(Default)).Methods("GET")
	router.HandleFunc("/"+ApiVersion+"/translations", wrap(Create)).Methods("POST")
	router.HandleFunc("/"+ApiVersion+"/translations/{id:[a-z0-9-]+}", wrap(Save)).Methods("POST")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}", wrap(Get)).Methods("GET")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}/{lang:[a-z]{2}}", wrap(GetParticular)).Methods("GET")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}", wrap(Delete)).Methods("DELETE")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}/{lang:[a-z]{2}}", wrap(DeleteParticular)).Methods("DELETE")

	app := negroni.New()
	//These middleware is common to all routes
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	//	app.Use(restgate.New(
	//		"X-Auth-Key",
	//		"X-Auth-Secret",
	//		restgate.Static,
	//		restgate.Config{
	//			Context: C,
	//			Key: []string{"12345"},
	//			Secret: []string{"secret"},
	//		},
	//	))
	app.UseHandler(router)
	http.Handle("/", context.ClearHandler(app))

	log.Printf("Info: Starting application on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	//	log.Fatal(http.ListenAndServe(":"+port, router))
	return nil
}

func cleanup() {
	driver.Client.Close()
	log.Print("Info: Gracefully closing application")
}

func C(r *http.Request, authenticatedKey string) {
	context.Set(r, 0, authenticatedKey)
}

func wrap(action Action) (func(http.ResponseWriter, *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request){
		v, code := action(w, r)

		callback := r.URL.Query().Get("callback")
		if callback == "" {
			render.JSON(w, code, v)
		} else {
			render.JSONP(w, code, callback, v)
		}
	}
}

func Ping(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	return "PONG", http.StatusOK
}

func Default(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	methodMap := make(map[string]string)
	methodMap["translation_map"] = fmt.Sprintf("/%s/%s", ApiVersion, "translations")
	return methodMap , http.StatusOK
}

func Create(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	u1 := uuid.NewV4().String()
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	bag, _ := SmartSave(request, u1)

	w.Header().Set("Location", "/"+ApiVersion+"/translations/"+u1)
	return bag, http.StatusCreated
}

func Save(w http.ResponseWriter, r *http.Request) (bag interface{}, status int) {
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	id := mux.Vars(r)["id"]
	bag, newLng := SmartSave(request, id)
	status = http.StatusOK

	if newLng > 0{
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

	id := mux.Vars(r)["id"]
	start := time.Now()
	bag, err := driver.GetAll(id)
	log.Printf("Completed in %v", time.Since(start))
	if nil != err {
		return map[string]string{"error":err.Error()}, http.StatusNotFound
	}

	return bag, http.StatusOK
}

func GetParticular(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := mux.Vars(r)["id"]
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
	id := mux.Vars(r)["id"]
	err := driver.Delete(id)
	if err != nil {
		log.Println("Delete", id, err)
	}
	return nil, http.StatusNoContent
}

func DeleteParticular(w http.ResponseWriter, r *http.Request) (interface{}, int) {
	id := mux.Vars(r)["id"]
	lang := mux.Vars(r)["lang"]
	err := driver.DeleteLang(id, lang)

	if err != nil {
		log.Println("Delete", id, err)
	}
	return nil, http.StatusNoContent
}

type RequestObject struct {
	Id     string `json:"id"`
	Text   string `json:"text"`
	Lang   []string `json:"lang"`
	Source string `json:"source"`
}


