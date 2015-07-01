package main

import (
	"log"
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

	log.Printf("Info: Starting application on port %s", port)

	log.Fatal(http.ListenAndServe(":8088", router))
	return nil
}

func cleanup() {
	log.Print("Info: Gracefully closing application")
}
