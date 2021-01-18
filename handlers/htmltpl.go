package handlers

import (
	"html/template"
	"io/ioutil"

	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
)

var tpl *template.Template
var tplLogin *template.Template

func init() {
	pkger.Include("/assets/src")

	{
		f, err := pkger.Open("/assets/src/login.html")
		if err != nil {
			logrus.Fatal("cannot open login template", err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			logrus.Fatal("cannot read login template", err)
		}

		tplLogin = template.Must(template.New("").Parse(string(b)))

	}
	{
		f, err := pkger.Open("/assets/src/page.html")
		if err != nil {
			logrus.Fatal("cannot open page template", err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			logrus.Fatal("cannot read page template", err)
		}

		tpl = template.Must(template.New("").Parse(string(b)))
	}

}
