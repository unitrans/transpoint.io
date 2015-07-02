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
)

const ApiVersion = "v1"

var (
	render *r.Render
)

func init() {
	godotenv.Load()
	render = r.New(r.Options{})
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

func cleanup() {
	log.Print("Info: Gracefully closing application")
}
