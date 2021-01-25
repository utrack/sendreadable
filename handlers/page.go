package handlers

import (
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

type pageRequest struct {
	Err            error
	LoggedIn       *pageLoginInfo
	DoLogout       bool
	customTpl      *template.Template
	SuccessMessage string
}

type pageLoginInfo struct {
	Email string
	Token string
}

func pageRenderErr(w http.ResponseWriter, r *http.Request, err error) {
	pageRender(w, r, pageRequest{Err: err})
}
func pageRender(w http.ResponseWriter, r *http.Request, rsp pageRequest) {
	w.WriteHeader(500)

	if rsp.LoggedIn != nil {
	}
	if rsp.DoLogout {
		coo := &http.Cookie{Name: cookieName,
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, coo)
	}

	w.Header().Set("Link", "</assets/style.css>; rel=preload;")

	ct := rsp.customTpl
	if ct == nil {
		ct = tpl
	}
	err := ct.Execute(w, rsp)
	if err != nil {
		logrus.Error("error when rendering the template", err)
	}
}
