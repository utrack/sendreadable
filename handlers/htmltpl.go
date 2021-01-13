package handlers

import (
	"crypto/sha1"
	"html/template"
	"io/ioutil"

	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
)

var tpl *template.Template
var tplEtag string

func init() {
	pkger.Include("/assets/src")

	f, err := pkger.Open("/assets/src/page.html")
	if err != nil {
		logrus.Fatal("cannot open page template", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.Fatal("cannot read page template", err)
	}

	tpl = template.Must(template.New("").Parse(string(b)))

	h := sha1.New()
	h.Write(b)
	tplEtag = string(h.Sum(nil))

}
