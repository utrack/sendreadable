package handlers

import (
	"net/http"

	"github.com/pkg/errors"
)

const cookieName = "tok"

func (h Handler) setCookie(w http.ResponseWriter, tok string) {

	coo := &http.Cookie{
		Name:     cookieName,
		Value:    tok,
		Secure:   h.secure,
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 7 * 4,
	}
	http.SetCookie(w, coo)
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		pageRender(w, r, pageRequest{
			customTpl: tplLogin,
		})
		return
	}

	r.ParseForm()

	code := r.FormValue("code")
	if code == "" {
		pageRender(w, r, pageRequest{Err: errors.New("Empty code provided"), customTpl: tplLogin})
		return
	}

	rsp, err := h.rm.Auth(r.Context(), code)
	if err != nil {
		pageRender(w, r, pageRequest{Err: errors.Wrap(err, "cannot get reMarkable user info"), customTpl: tplLogin})
		return
	}

	tok, err := jwtGen(h.jwtKey, rsp.Tokens)
	if err != nil {
		pageRender(w, r, pageRequest{Err: errors.Wrap(err, "cannot generate JWT token"), customTpl: tplLogin})
	}

	// if strings.HasPrefix(r.Header.Get("referer"), "http") {
	// 	isSecure = false
	// }

	h.setCookie(w, tok)

	http.Redirect(w, r, "/", 303)
}

func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	coo := &http.Cookie{Name: cookieName,
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, coo)

	http.Redirect(w, r, "/", 303)
}
