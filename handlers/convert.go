package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/utrack/sendreadable/converter"
)

type Handler struct {
	svc           *converter.Service
	chopDirPrefix string
}

func New(svc *converter.Service, dirPrefix string) Handler {
	return Handler{svc: svc, chopDirPrefix: dirPrefix}
}

func (h Handler) Convert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := r.FormValue("url")
	if u == "" {
		w.Write([]byte(page))
		return
	}
	res, err := h.svc.Convert(r.Context(), u)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	newPath := "/download/" + strings.TrimSuffix(strings.TrimPrefix(res.Filename, h.chopDirPrefix), "/main.pdf")
	artName := url.QueryEscape(res.ArticleName)

	q := newPath + "?filename=" + artName

	http.Redirect(w, r, q, 307)
}
