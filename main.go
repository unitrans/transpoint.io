package main

import (
	"os"
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
	port := os.Getenv("PORT")
	router := mux.NewRouter()
	router.HandleFunc("/ping", Ping).Methods("GET")

	log.Printf("Info: Starting application on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
	return nil
}

func Ping (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PONG")
}

func cleanup() {
	log.Print("Info: Gracefully closing application")
}
