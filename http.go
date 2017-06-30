package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type myHandler struct {
	cfg   oauth2.Config
	state string
	token *oauth2.Token
}

func newHandler(cfg oauth2.Config) http.Handler {
	state := randStringRunes(24)

	h := myHandler{
		cfg:   cfg,
		state: state,
	}

	return h
}

func registerServeMux(handler http.Handler) {
	http.Handle("/", handler)
}

func (m myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/oauth2/callback" {
		if r.URL.Query().Get("state") != m.state {
			http.Error(w, "bad state", http.StatusForbidden)
			return
		}
		fmt.Println("Got callback: %+v\n", r.URL.Query())
		return
	}

	if m.token == nil {
		redir := m.cfg.AuthCodeURL(m.state)
		http.Redirect(w, r, redir, http.StatusFound)
		return
	}

	if m.token.Valid() {
		fmt.Fprintf(w, "%+v\n", m.token)
	}

	// TODO
	io.WriteString(w, "ok\n")
}
