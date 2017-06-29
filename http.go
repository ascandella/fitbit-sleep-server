package main

import (
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type myHandler struct {
	cfg   oauth2.Config
	state string
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
	if r.URL.Path == "/oauth2/authorize" {
		// TODO
	}
	// TODO
	io.WriteString(w, "ok\n")
}
