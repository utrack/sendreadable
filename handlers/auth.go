package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/utrack/sendreadable/pkg/rmclient"
)

type Auth struct {
	rm *rmclient.Client

	jwtKey interface{}
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	code := r.FormValue("code")
	if code == "" {
		pageRenderErr(w, r, errors.New("Empty code provided"))
		return
	}

	rsp, err := a.rm.Auth(r.Context(), code)
	if err != nil {
		pageRenderErr(w, r, errors.Wrap(err, "cannot get reMarkable user info"))
		return
	}

	tok, err := jwtGen(a.jwtKey, rsp.Token)
	if err != nil {
		pageRenderErr(w, r, errors.Wrap(err, "cannot generate JWT token"))
	}

	pageRender(w, r, pageRequest{
		LoggedIn: &pageLoginInfo{
			Token: tok,
		},
	})

}
