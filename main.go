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

	"encoding/json"
)

const ApiVersion = "v1"

func init() {
	godotenv.Load()
}

func main() {
	closer.Bind(cleanup)
	closer.Checked(run, true)
}

func run() error {
	port := os.Getenv("PORT")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", Ping).Methods("GET")
	router.HandleFunc("/"+ApiVersion+"/", Default).Methods("GET")

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

func C(r *http.Request, authenticatedKey string) {
	context.Set(r, 0, authenticatedKey) // Read http://www.gorillatoolkit.org/pkg/context about setting arbitary context key
}

func ResponseJson(w http.ResponseWriter, v interface{}) {

	jsonBytes, err := json.Marshal(v)
	if err != nil{
		jsonBytes, _ = json.Marshal(fmt.Sprintf("%v", v))
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(jsonBytes))
}

func Ping (w http.ResponseWriter, r *http.Request) {
	ResponseJson(w, "PONG")
}

func Default (w http.ResponseWriter, r *http.Request) {
	methodMap := make(map[string]string)
	methodMap["translation_map"] = fmt.Sprintf("/%s/%s", ApiVersion, "translations")
	ResponseJson(w, methodMap)
}

func cleanup() {
	log.Print("Info: Gracefully closing application")
}
