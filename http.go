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

	"go.uber.org/zap"

	"golang.org/x/oauth2"
)

type myHandler struct {
	cfg       oauth2.Config
	state     string
	token     *oauth2.Token
	tknSource oauth2.TokenSource
	client    *http.Client
	log       *zap.Logger
}

func newHandler(cfg oauth2.Config, log *zap.Logger) *myHandler {
	state := randStringRunes(24)

	h := myHandler{
		cfg:   cfg,
		state: state,
		log:   log,
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
		m.log.Info("Redirecting user to oauth2 fitbit", zap.String("url", redir))
		http.Redirect(w, r, redir, http.StatusFound)
	case "/":
		if m.token.Valid() {
			m.log.Info("Have valid token")
			m.getAndCacheSleep(w, r)
			return
		}
		m.log.Error("No token -- need to hit /aiden to authorize")
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
		m.log.Error("Got error from from Fitbit API", zap.String("url", u), zap.Error(err))
		return
	}

	if sleep.StatusCode != http.StatusOK {
		http.Error(w, "Got non-200 from fitbit sleep API", http.StatusInternalServerError)
		fmt.Printf("Got bad status from %q: %s\n", u, sleep.Status)
		io.Copy(os.Stdout, sleep.Body)
		return
	}

	m.log.Info("Got 200 OK from fitbit API", zap.String("url", u))

	defer func() {
		if err := sleep.Body.Close(); err != nil {
			m.log.Error("Couldn't close body", zap.Error(err))
		}
	}()

	log := sleepLog{}
	dec := json.NewDecoder(sleep.Body)
	if err := dec.Decode(&log); err != nil {
		fmt.Println("Error decoding body: ", err.Error())
		return
	}

	m.log.Info("Received data", zap.Any("data", log))
	sleepTemplate.Execute(w, log.Sleep[0])
}

func (m *myHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != m.state {
		http.Error(w, "bad state", http.StatusForbidden)
		return
	}
	m.log.Info("Got callback: %+v\n", zap.String("query", r.URL.Query().Encode()))

	tkn, err := m.cfg.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		m.log.Error("Error exchanging code", zap.Error(err))
		http.Error(w, "Unable to exchange oauth2 code", http.StatusInternalServerError)
		return
	}
	m.registerToken(tkn)

	fmt.Println("Redirecting home")
	http.Redirect(w, r, "/", http.StatusFound)
}

func (m *myHandler) registerToken(tkn *oauth2.Token) {
	m.token = tkn

	m.log.Info("Registering token", zap.String("token", tkn.AccessToken))
	bs, _ := json.Marshal(tkn)
	fmt.Printf(string(bs))
	m.tknSource = oauth2.ReuseTokenSource(tkn, nil)

	m.log.Debug("Got token source", zap.String("source", fmt.Sprintf("%s", m.tknSource)))

	m.client = oauth2.NewClient(context.Background(), m.tknSource)
}
