package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	_ "github.com/utrack/sendreadable/assets/compiled"
	"github.com/utrack/sendreadable/pkg/rmclient"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/markbates/pkger"
	"github.com/utrack/sendreadable/converter"
	"github.com/utrack/sendreadable/handlers"
)

func init() {
	pkger.Include("/assets/dist")
	os.Stderr = os.Stdout
}

func main() {
	defer os.Stdout.Sync()

	pathToJwtKey := flag.String("key", "priv.key", "path to JWT private key")
	secureCookies := flag.Bool("secure", false, "send ssl-only cookies")
	flag.Parse()

	keyBuf, err := ioutil.ReadFile(*pathToJwtKey)
	if err != nil {
		log.Fatal("reading JWT key: ", err.Error())
	}
	key, err := parseJWTKey(keyBuf)
	if err != nil {
		log.Fatal("parsing JWT key: ", err.Error())
	}

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

	hdl := handlers.New(svc, tmpDir, rmclient.New(), key, *secureCookies)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.HandleFunc("/conv/conv/assets/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=600, stale-if-error=3600")
		w.Header().Set("Content-Type", "text/css")
		w.Write(style)
	})
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
	r.Get("/", http.HandlerFunc(hdl.Convert))
	r.Post("/", http.HandlerFunc(hdl.Convert))
	r.Get("/login", http.HandlerFunc(hdl.Login))
	r.Post("/login", http.HandlerFunc(hdl.Login))
	r.Get("/logout", http.HandlerFunc(hdl.Logout))

	http.ListenAndServe(":3333", r)
}

func parseJWTKey(b []byte) (*rsa.PrivateKey, error) {
	privPem, _ := pem.Decode(b)

	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.Errorf("wrong RSA key type: %v", privPem.Type)
	}

	privBytes := privPem.Bytes

	var parsedKey interface{}
	var err error
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privBytes); err != nil {
			return nil, err
		}
	}

	privateKey := parsedKey.(*rsa.PrivateKey)
	return privateKey, nil
}
