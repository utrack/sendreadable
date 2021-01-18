package handlers

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (h Handler) Convert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := r.FormValue("url")

	sendRemarkable := false
	if r.Method == "POST" {
		sendRemarkable = true
	}

	var preq pageRequest
	var rmToken string

	tok, err := r.Cookie(cookieName)
	if err == nil {
		var cl jwtClaims
		_, err := jwt.ParseWithClaims(tok.Value, &cl,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != jwt.SigningMethodRS512.Alg() {
					return nil, errors.Errorf("bad alg '%v'", token.Method.Alg())
				}
				return &h.jwtKey.PublicKey, nil
			})
		if err != nil || cl.Valid() != nil {
			logrus.Warn("bad JWT token, logging out: ", err, ", ", cl.Valid())
			preq.DoLogout = true
		} else {
			preq.LoggedIn = &pageLoginInfo{}
			rmToken = cl.RmJWT
		}
	}
	if sendRemarkable && rmToken == "" {
		preq.DoLogout = true
		preq.Err = errors.New("cannot send to reMarkable: not logged in or token is stale")
		pageRender(w, r, preq)
		return
	}

	if u == "" {
		pageRender(w, r, preq)
		return
	}
	res, err := h.svc.Convert(r.Context(), u)
	if err != nil {
		w.WriteHeader(500)
		pageRenderErr(w, r, err)
		return
	}

	if !sendRemarkable {
		newPath := "/download/" + strings.TrimSuffix(strings.TrimPrefix(res.Filename, h.chopDirPrefix), "/main.pdf")
		artName := url.QueryEscape(res.ArticleName)
		q := newPath + "?filename=" + artName

		http.Redirect(w, r, q, 307)
		return
	}

	f, err := os.Open(res.Filename)
	if err != nil {
		preq.Err = errors.Wrap(err, "cannot open tmpfile")
		pageRender(w, r, preq)
		return
	}
	defer f.Close()

	err = h.rm.Upload(r.Context(), res.ArticleName+".pdf", f, rmToken)
	preq.Err = errors.Wrap(err, "cannot send file to reMarkable")
	if err == nil {
		preq.SuccessMessage = "Sent '" + res.ArticleName + ".pdf'"
	}
	pageRender(w, r, preq)
}
