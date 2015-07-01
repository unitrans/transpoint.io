// Copyright ${YEAR} Home24 AG. All rights reserved.
// Proprietary license.

package main
import "github.com/gorilla/mux"

func Router() *mux.Router {

	//Create subrouters
//	restRouter := mux.NewRouter()
//	restRouter.HandleFunc("/api", Handler1()) //Rest API Endpoint handler -> Use your own
//
//	rest2Router := mux.NewRouter()
//	rest2Router.HandleFunc("/api2", Handler2()) //A second Rest API Endpoint handler -> Use your own
//
//	//Create negroni instance to handle different middlewares for different api routes
//	negRest := negroni.New()
//	negRest.Use(restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Static, restgate.Config{Context: C, Key: []string{"12345"}, Secret: []string{"secret"}}))
//	negRest.UseHandler(restRouter)
//
//	negRest2 := negroni.New()
//	negRest2.Use(restgate.New("X-Auth-Key", "X-Auth-Secret", restgate.Database, restgate.Config{DB: SqlDB(), TableName: "users", Key: []string{"keys"}, Secret: []string{"secrets"}}))
//	negRest2.UseHandler(rest2Router)
//
//	//Create main router
//	mainRouter := mux.NewRouter().StrictSlash(true)
//	mainRouter.HandleFunc("/", MainHandler()) //Main Handler -> Use your own
//	mainRouter.Handle("/api", negRest) //This endpoint is protected by RestGate via hardcoded KEYs
//	mainRouter.Handle("/api2", negRest2) //This endpoint is protected by RestGate via KEYs stored in a database
//
//	return mainRouter

}