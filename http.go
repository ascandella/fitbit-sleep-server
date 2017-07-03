package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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
	switch r.URL.Path {
	case "/oauth2/callback":
		m.handleCallback(w, r)
	case "/aiden":
		redir := m.cfg.AuthCodeURL(m.state)
		fmt.Println("Redirecting user to oauth2 fitbit: ", redir)
		http.Redirect(w, r, redir, http.StatusFound)
	case "/":
		if m.token.Valid() {
			fmt.Println("Have valid token")
			m.getAndCacheSleep(w, r)
			return
		}

		fmt.Println("No token -- need to hit /aiden to authorize")
		http.Error(w, "no token -- aiden must have screwed something up", http.StatusInternalServerError)
	}
}

const sleepEndpoint = "https://api.fitbit.com/1.2/user/-/sleep/"

func (m *myHandler) getAndCacheSleep(w http.ResponseWriter, r *http.Request) {
	// TODO caching
	u := sleepEndpoint
	date := r.URL.Query().Get("date")
	if date != "" {
		u += "date/" + date + ".json"
	} else {
		afterTime := time.Now().Add(-72 * time.Hour)
		after := strings.Split(afterTime.Format(time.RFC3339), "T")[0]
		u += "list.json?limit=3&offset=0&sort=desc&afterDate=" + after
	}

	sleep, err := m.client.Get(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Got error from %q: %s\n", u, err.Error())
		return
	}

	if sleep.StatusCode != http.StatusOK {
		http.Error(w, "Got non-200 from fitbit sleep API", http.StatusInternalServerError)
		fmt.Printf("Got bad status from %q: %s\n", u, sleep.Status)
		io.Copy(os.Stdout, sleep.Body)
		return
	}

	fmt.Printf("Got 200 OK from fitbit API for %q\n", u)

	defer func() {
		if err := sleep.Body.Close(); err != nil {
			fmt.Printf("[ERROR] Couldn't close body: %s", err.Error())
		}
	}()

	log := sleepLog{}
	dec := json.NewDecoder(sleep.Body)
	if err := dec.Decode(&log); err != nil {
		fmt.Println("Error decoding body: ", err.Error())
		return
	}

	fmt.Printf("Data: %+v\n", log)
	sleepTemplate.Execute(w, log.MostRecent())
}

func (m *myHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
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
	m.registerToken(tkn)

	fmt.Println("Redirecting home")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (m *myHandler) registerToken(tkn *oauth2.Token) {
	m.token = tkn

	fmt.Printf("Token: %+v\n", tkn)
	m.tknSource = oauth2.ReuseTokenSource(tkn, nil)

	fmt.Println("Got token source: ", m.tknSource)

	m.client = oauth2.NewClient(context.Background(), m.tknSource)

}
