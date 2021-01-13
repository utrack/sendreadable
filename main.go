package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/utrack/sendreadable/assets/compiled"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/markbates/pkger"
	"github.com/utrack/sendreadable/converter"
	"github.com/utrack/sendreadable/handlers"
)

func main() {
	pkger.Include("/assets/dist")

	st, err := pkger.Open("/assets/dist/style.css")
	if err != nil {
		log.Fatal("opening style:", err.Error())
	}
	style, err := ioutil.ReadAll(st)
	if err != nil {
		log.Fatal(err.Error())
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err.Error())
	}

	tmpDir := os.TempDir() + "/sendreadable"

	svc := converter.New(dir+"/fonts", tmpDir)

	hdl := handlers.New(svc, tmpDir)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.HandleFunc("/conv/assets/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=600, stale-if-error=3600")
		w.Header().Set("Content-Type", "text/css")
		w.Write(style)
	})
	r.HandleFunc("/assets/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=600, stale-if-error=3600")
		w.Header().Set("Content-Type", "text/css")
		w.Write(style)
	})

	r.Get("/conv", http.HandlerFunc(hdl.Convert))

	http.ListenAndServe(":3333", r)
}
