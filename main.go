package main

import (
	"os"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
	"github.com/davecheney/profile"
//	"github.com/xlab/closer"
	r "gopkg.in/unrolled/render.v1"

	"github.com/urakozz/transpoint.io/storage"
	t "github.com/urakozz/transpoint.io/translator"
	"github.com/urakozz/transpoint.io/src/infrastrucrute"
)


var (
	render *r.Render
	driver *storage.RedisDriver
	translator *t.YandexTranslator
	cookieStore *sessions.CookieStore
	profiler interface { Stop() }
)

type Action func(w http.ResponseWriter, r *http.Request) (interface{}, int)

func init() {
	godotenv.Load()
	render = r.New(r.Options{})
	driver = storage.NewRedisDriver(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"))
	translator = t.NewYandexTranslator(os.Getenv("Y_TR_KEY"))
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
	//	closer.Bind(cleanup)
	//	closer.Checked(run, true)
	run()
}

func run() error {
	port := os.Getenv("PORT")

	http.Handle("/v1/", http.StripPrefix("/v1/", context.ClearHandler(ApiRouter())))
	http.HandleFunc("/ping", ApiPing())
	http.Handle("/webapp/", http.StripPrefix("/webapp", context.ClearHandler(WebRouter())))
	http.Handle("/webapi/", http.StripPrefix("/webapi", infrastrucrute.NewWebApi()))
	http.HandleFunc("/", WebIndexPage)

	initProfiler()

	log.Printf("Info: Starting application on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	//	log.Fatal(http.ListenAndServe(":"+port, router))
	return nil
}

func cleanup() {
	driver.Client.Close()
	log.Print("Info: Gracefully closing application")
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



