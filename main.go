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
	"github.com/xlab/closer"
	"github.com/pjebs/restgate"
	r "gopkg.in/unrolled/render.v1"

	"github.com/urakozz/transpoint.io/storage"
	"github.com/satori/go.uuid"
	"encoding/json"
	"time"
	"github.com/urakozz/transpoint.io/translator"
)

const ApiVersion = "v1"

var (
	render *r.Render
	driver *storage.RedisDriver
)

func init() {
	godotenv.Load()
	render = r.New(r.Options{})
	driver = storage.NewRedisDriver(os.Getenv("REDIS_ADDR"))

}

func main() {
	closer.Bind(cleanup)
	closer.Checked(run, true)
}

func run() error {
	port := os.Getenv("PORT")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", Ping).Methods("GET")
	router.HandleFunc("/"+ApiVersion, Default).Methods("GET")
	router.HandleFunc("/"+ApiVersion+"/translations", Create).Methods("POST")
	router.HandleFunc("/"+ApiVersion+"/translations/{id:[a-z0-9-]+}", Save).Methods("POST")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}", Get).Methods("GET")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}/{lang:[a-z]{2}}", GetParticular).Methods("GET")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}", Delete).Methods("DELETE")
	router.HandleFunc("/"+ApiVersion+"/translations/{id}/{lang:[a-z]{2}}", DeleteParticular).Methods("DELETE")

	app := negroni.New()
	//These middleware is common to all routes
	app.Use(negroni.NewRecovery())
	app.Use(negroni.NewLogger())
	app.Use(restgate.New(
		"X-Auth-Key",
		"X-Auth-Secret",
		restgate.Static,
		restgate.Config{
			Context: C,
			Key: []string{"12345"},
			Secret: []string{"secret"},
		},
	))
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

func Ping (w http.ResponseWriter, r *http.Request) {
	render.JSON(w, http.StatusOK, "PONG")
}

func Default (w http.ResponseWriter, r *http.Request) {
	methodMap := make(map[string]string)
	methodMap["translation_map"] = fmt.Sprintf("/%s/%s", ApiVersion, "translations")
	render.JSON(w, http.StatusOK, methodMap)
}

func Create (w http.ResponseWriter, r *http.Request) {
	u1 := uuid.NewV4().String()
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	translator := translator.NewYandexTranslator()
	container := translator.Translate(request.Text, request.Lang)
	log.Println(container)
	go func(){
		driver.Save(u1, container.Source, request.Text, container.Bag)
	}()
	w.Header().Set("Location", "/"+ApiVersion+"/translations/"+u1)
	render.JSON(w, http.StatusCreated, u1)
}

func Save (w http.ResponseWriter, r *http.Request) {
	var request *RequestObject
	json.NewDecoder(r.Body).Decode(&request)
	id := mux.Vars(r)["id"]

	go func(){
		driver.Save(id, "ru", request.Text, map[string]string{"ru": "rrrr", "en": "eeeee"})
	}()
	w.Header().Set("Location", "/"+ApiVersion+"/translations/"+id)
	render.JSON(w, http.StatusCreated, "Save")
}

func Get (w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	start := time.Now()
	bag, err := driver.GetAll(id)
	log.Printf("Completed in %v", time.Since(start))
	if nil != err {
		render.JSON(w, http.StatusNotFound, map[string]string{"error":err.Error()})
		return
	}

	render.JSON(w, http.StatusOK, bag)
}

func GetParticular (w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	lang := mux.Vars(r)["lang"]

	start := time.Now()
	bag, err := driver.GetLang(id, lang)
	log.Printf("Completed in %v", time.Since(start))

	if nil != err {
		render.JSON(w, http.StatusNotFound, map[string]string{"error":err.Error()})
		return
	}

	render.JSON(w, http.StatusOK, bag)
}

func Delete (w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := driver.Delete(id)
	if err != nil {
		log.Println("Delete", id, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteParticular (w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	lang := mux.Vars(r)["lang"]
	err := driver.DeleteLang(id, lang)

	if err != nil {
		log.Println("Delete", id, err)
	}
	w.WriteHeader(http.StatusNoContent)
}

type RequestObject struct {
	Id string `json:"id"`
	Text string `json:"text"`
	Lang []string `json:"lang"`
	Source string `json:"source"`
}


