package main

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

const secretsJSON = "secrets.json"

func loadConfigFromJSON(location string) (oauth2.Config, error) {
	c := oauth2.Config{
		Endpoint: fitbit.Endpoint,
	}

	f, err := os.Open(location)
	if err != nil {
		return c, err
	}

	p := json.NewDecoder(f)
	if err := p.Decode(&c); err != nil {
		return c, err
	}

	return c, nil
}
