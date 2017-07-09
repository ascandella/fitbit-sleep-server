package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"golang.org/x/oauth2"
)

const (
	sleepEndpoint = "https://api.fitbit.com/1.2/user/-/sleep/"
	cacheTime     = 30 * time.Minute
)

type myHandler struct {
	cfg       oauth2.Config
	state     string
	token     *oauth2.Token
	log       *zap.Logger
	appConfig appConfig
	mu        sync.Mutex
	lastFetch time.Time
	cachedLog sleepLog
}

func newHandler(cfg oauth2.Config, appConf appConfig) *myHandler {
	state := randStringRunes(24)

	h := myHandler{
		appConfig: appConf,
		cfg:       cfg,
		log:       appConf.log,
		mu:        sync.Mutex{},
		state:     state,
		cachedLog: sleepLog{},
		lastFetch: time.Now().Add(-cacheTime),
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
		redir := m.cfg.AuthCodeURL(m.state, oauth2.AccessTypeOffline)
		m.log.Info("Redirecting user to oauth2 fitbit", zap.String("url", redir))
		http.Redirect(w, r, redir, http.StatusFound)
	case "/":
		m.log.Info("Serving home - token valid?", zap.Bool("valid", m.token.Valid()))
		if !m.token.Valid() {
			m.log.Warn("Will need token refresh")
		}
		m.getAndCacheSleep(w, r)
	}
}

func (m *myHandler) getAndCacheSleep(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delta := time.Now().Sub(m.lastFetch)
	if delta < cacheTime && len(m.cachedLog.Sleep) > 0 {
		m.log.Info("Returning cached version", zap.Duration("delta", delta))
		m.showLog(w, m.cachedLog)
		return
	}

	afterTime := time.Now().Add(-72 * time.Hour)
	after := strings.Split(afterTime.Format(time.RFC3339), "T")[0]
	u := sleepEndpoint + "list.json?limit=3&offset=0&sort=desc&afterDate=" + after

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := m.cfg.Client(ctx, m.token)
	sleep, err := client.Get(u)
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

	m.lastFetch = time.Now()
	m.cachedLog = log

	m.showLog(w, log)
}

func (m *myHandler) showLog(w http.ResponseWriter, log sleepLog) {
	m.log.Info("Displaying data", zap.Any("data", log))
	sleepTemplate.Execute(w, log.Sleep[0])
}

func (m *myHandler) handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != m.state {
		m.log.Error("State mismatch", zap.String("ours", m.state), zap.String("theirs", r.URL.Query().Get("state")))
		http.Error(w, "bad state", http.StatusForbidden)
		return
	}
	m.log.Info("Handling oauth2 callback with good state", zap.String("query", r.URL.Query().Encode()))

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

func (m *myHandler) maybeStoreToken(tkn *oauth2.Token) {
	if m.appConfig.tokenFile == nil || *m.appConfig.tokenFile == "" {
		m.log.Debug("No token file specified, not persisting.")
		return
	}

	bs, err := json.Marshal(tkn)
	if err != nil {
		m.log.Error("Unable to marshal token", zap.Error(err))
		return
	}

	m.log.Info("Writing token to file", zap.String("path", *m.appConfig.tokenFile))
	if err := ioutil.WriteFile(*m.appConfig.tokenFile, bs, 0600); err != nil {
		m.log.Error("Unable to save token file", zap.Error(err))
	}
}

func (m *myHandler) registerToken(tkn *oauth2.Token) {
	m.log.Info("Registering token", zap.Any("token", tkn))
	m.token = tkn

	m.maybeStoreToken(tkn)
}
