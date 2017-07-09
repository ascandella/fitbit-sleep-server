package main

import (
	"testing"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
)

func TestNewHandler_OK(t *testing.T) {
	h := newHandler(oauth2.Config{}, appConfig{})
	assert.NotEmpty(t, h.state)
}

func TestRegisterServMux(t *testing.T) {
	registerServeMux(newHandler(oauth2.Config{}, appConfig{}))
}
