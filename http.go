package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type myHandler struct {
	cfg       oauth2.Config
	state     string
	token     *oauth2.Token
	tknSource oauth2.TokenSource
	client    *http.Client
}

func newHandler(cfg oauth2.Config) http.Handler {
	state := randStringRunes(24)

	h := myHandler{
		cfg:   cfg,
		state: state,
	}

	return &h
}

func registerServeMux(handler http.Handler) {
	http.Handle("/", handler)
}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/oauth2/callback" {
		if r.URL.Query().Get("state") != m.state {
			http.Error(w, "bad state", http.StatusForbidden)
			return
		}
		fmt.Printf("Got callback: %+v\n", r.URL.Query())

		tkn, err := m.cfg.Exchange(context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			fmt.Printf("Error exchanging code: %+v\n", err)
			http.Error(w, "Unable to exchange oauth2 code", http.StatusInternalServerError)
			return
		}
		m.token = tkn

		fmt.Printf("Token: %+v\n", tkn)
		m.tknSource = oauth2.ReuseTokenSource(tkn, nil)

		fmt.Println("Got token source: ", m.tknSource)

		m.client = oauth2.NewClient(context.Background(), m.tknSource)

		fmt.Println("Redirecting home")
		http.Redirect(w, r, "/", http.StatusFound)

		return
	}

	if m.token == nil {
		redir := m.cfg.AuthCodeURL(m.state)
		fmt.Println("Redirecting user to oauth2 fitbit: ", redir)
		http.Redirect(w, r, redir, http.StatusFound)
		return
	}

	if m.token.Valid() {
		fmt.Fprintf(w, "Have valid token!")
		return
	}

	// TODO
	io.WriteString(w, "shouldn't get here\n")
}
