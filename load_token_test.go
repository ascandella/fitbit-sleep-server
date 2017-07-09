package main

import (
	"testing"

	"code.ndella.com/ai-life/internal/testutils"
	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadToken_NoFile(t *testing.T) {
	_, err := loadTokenFromFile("non-existent")
	assert.Error(t, err)
}

func TestLoadToken_ValidVile(t *testing.T) {
	secret := "sooper-secret"
	tkn := oauth2.Token{AccessToken: secret}
	testutils.WithJSONFile(t, tkn, func(path string) {
		token, err := loadTokenFromFile(path)
		require.NoError(t, err)
		require.Equal(t, token.AccessToken, secret)
	})
}
