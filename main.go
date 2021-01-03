package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/utrack/sendreadable/converter"
	"github.com/utrack/sendreadable/handlers"
)

func main() {

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

	r.Get("/conv", http.HandlerFunc(hdl.Convert))

	http.ListenAndServe(":3333", r)
}
