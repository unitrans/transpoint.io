package main

import (
	"os"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/davecheney/profile"
	"gopkg.in/redis.v3"

	t "github.com/urakozz/transpoint.io/src/translator"
	"github.com/urakozz/transpoint.io/src/scenario"

	"github.com/urakozz/transpoint.io/src/infrastrucrute/storage"
	"github.com/urakozz/transpoint.io/src/interface/repository/redis"
	//"github.com/urakozz/transpoint.io/src/infrastrucrute/translator/particular"
	"github.com/urakozz/transpoint.io/src/translator/backend_full"
	"github.com/urakozz/transpoint.io/src/infrastrucrute/httpclient"
)


var (
	redisClient *redis.Client
	translator t.Translator
	cookieStore *sessions.CookieStore
	profiler interface { Stop() }
)


func init() {
	godotenv.Load()
	redisClient = storage.RedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"))
	translator = t.NewTranslateAdapter(
		[]backend_full.IBackendFull{
			backend_full.NewGoogleTranslator(httpclient.GetHttpClient(), os.Getenv("G_TR_KEY")),
			backend_full.NewYandexTranslator(httpclient.GetHttpClient(), os.Getenv("Y_TR_KEY")),
//			backend_full.NewBingTranslator(os.Getenv("B_TR_KEY")),
		})
	//translator.AddParticular(&particular.AbbyyLingvoLiveTranslator{})
	if "" == os.Getenv("APP_SECRET") {
		os.Setenv("APP_SECRET", string(securecookie.GenerateRandomKey(32)))
	}
	cookieStore = &sessions.CookieStore{
		Codecs: securecookie.CodecsFromPairs([]byte(os.Getenv("APP_SECRET"))),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30 * 10,
			//			Secure:true,
			HttpOnly: true,
		},
	}
}

func main() {
	port := os.Getenv("PORT")


	userRepository := repository.NewUserRepository(redisClient)
	transRepository := repository.NewTranslationRepository(redisClient)

	http.Handle("/v1/", http.StripPrefix("/v1", scenario.ApiRouter(userRepository, transRepository, translator)))
	http.Handle("/webapp/", http.StripPrefix("/webapp", scenario.WebRouter(userRepository, cookieStore)))
	http.Handle("/webapi/", http.StripPrefix("/webapi", scenario.NewWebApi()))
	http.HandleFunc("/ping", scenario.ApiPing())
	http.HandleFunc("/", scenario.WebIndexPage)

	//initProfiler()

	log.Printf("Info: Starting application on port %s", port)
	log.Printf("Info: app done")
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

func initProfiler() {
	cfg := profile.Config{
		MemProfile:     true,
		CPUProfile:     true,
		ProfilePath:    ".", // store profiles in current directory
	}

	// p.Stop() must be called before the program exits to
	// ensure profiling information is written to disk.
	profiler = profile.Start(&cfg)
}



