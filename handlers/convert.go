package handlers

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/utrack/sendreadable/pkg/rmclient"
)

func (h Handler) Convert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := r.FormValue("url")

	sendRemarkable := false
	if r.Method == "POST" {
		sendRemarkable = true
	}

	var preq pageRequest
	var rmToken rmclient.Tokens

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
			if err != nil {
				preq.Err = errors.Wrap(err, "cannot parse rM JWT token")
			} else {
				preq.Err = errors.New("rM JWT token is invalid")
			}
			preq.Err = errors.Wrap(err, "cannot send to reMarkable")
			preq.DoLogout = true
			pageRender(w, r, preq)
			return
		} else {
			preq.LoggedIn = &pageLoginInfo{}
			rmToken.User = cl.RmUserTok
			rmToken.Device = cl.RmDeviceTok
			rmToken.UserRefreshedAt = time.Unix(cl.IssuedAt, 0)
		}
	}
	if sendRemarkable && rmToken.Device == "" && rmToken.User == "" {
		preq.DoLogout = true
		preq.Err = errors.New("empty token")
		pageRender(w, r, preq)
	}

	if u == "" {
		pageRender(w, r, preq)
		return
	}
	res, err := h.svc.Convert(r.Context(), u)
	if err != nil {
		preq.Err = err
		pageRender(w, r, preq)
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

	oldTok := rmToken
	err = h.rm.Upload(r.Context(), res.ArticleName, f, &rmToken)
	preq.Err = errors.Wrap(err, "cannot send file to reMarkable")
	if err == nil {
		preq.SuccessMessage = "Sent '" + res.ArticleName + "'"
	}
	if oldTok != rmToken {
		tok, err := jwtGen(h.jwtKey, rmToken)
		if err == nil {
			h.setCookie(w, tok)
		}
	}
	pageRender(w, r, preq)
}
