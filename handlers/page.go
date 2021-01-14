package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type pageRequest struct {
	Err      error
	LoggedIn *pageLoginInfo
}

type pageLoginInfo struct {
	Email string
	Token string
}

func pageRenderErr(w http.ResponseWriter, r *http.Request, err error) {
	pageRender(w, r, pageRequest{Err: err})
}
func pageRender(w http.ResponseWriter, r *http.Request, rsp pageRequest) {
	w.Header().Set("Link", "</assets/style.css>; rel=preload;")
	if rsp.Err == nil {
		w.Header().Set("Cache-Control", "public, max-age=3600, stale-if-error=60")
	} else {
		w.Header().Set("Cache-Control", "public, max-age=60, stale-if-error=60")
	}
	err := tpl.Execute(w, rsp)
	if err != nil {
		logrus.Error("error when rendering the template", err)
	}
}
