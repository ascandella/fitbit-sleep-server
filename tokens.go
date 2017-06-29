package main

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

const secretsJSON = "secrets.json"

func getToken() oauth2.TokenSource {

	return oauth2.ReuseTokenSource(loadTokensFromJSON(secretsJSON), nil)
}

func loadConfigFromJSON(location string) (oauth2.Config, error) {
	c := oauth2.Config{
		Endpoint: fitbit.Endpoint,
	}

	f, err := os.Open(location)
}
