package main

import (
	"net/http/httptest"
	"os"
	"testing"

	"go.uber.org/zap"

	"golang.org/x/oauth2"

	"github.com/sectioneight/fitbit-sleep-server/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler_OK(t *testing.T) {
	h := emptyHandler()
	assert.NotEmpty(t, h.state)
}

func TestRegisterServMux(t *testing.T) {
	registerServeMux(emptyHandler())
}

func TestRegisterToken_OK(t *testing.T) {
	h := emptyHandler()

	tkn := &oauth2.Token{
		AccessToken: "foo",
	}

	testutils.WithTempDir(t, func(dir string) {
		testutils.WithFile(t, dir, func(file *os.File) {
			fname := file.Name()
			h.appConfig.tokenFile = &fname
			h.registerToken(tkn)
		})
	})
}

func TestMaybeStoreToken_NoFile(t *testing.T) {
	h := emptyHandler()
	h.maybeStoreToken(nil)
}

func TestShowLog_NoData(t *testing.T) {
	w := httptest.NewRecorder()
	h := emptyHandler()
	h.showLog(w, sleepLog{})
	assert.Contains(t, w.Body.String(), "sleep data")
}

func TestShowLog_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	h := emptyHandler()
	r := httptest.NewRequest("GET", "/aoeu", nil)
	h.ServeHTTP(w, r)
	body := w.Body.String()
	assert.Contains(t, body, "no such page")
	assert.Contains(t, body, "aoeu")
}

func emptyHandler() *myHandler {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return newHandler(oauth2.Config{}, appConfig{
		log: l,
	})
}
