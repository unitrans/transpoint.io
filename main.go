package main

import (
	"log"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/xlab/closer"

)

func init() {
	godotenv.Load()
}

func main() {
	closer.Bind(cleanup)
	closer.Checked(run, true)
}

func run() error {

	router := mux.NewRouter()
	router.HandleFunc("/ping", Ping).Methods("GET")

	log.Printf("Info: Starting application on port %s", "8088")

	log.Fatal(http.ListenAndServe(":8088", router))
	return nil
}

func Ping (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PONG")
}

func cleanup() {
	log.Print("Info: Gracefully closing application")
}
