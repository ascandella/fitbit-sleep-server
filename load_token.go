package main

import (
	"encoding/json"
	"os"

	"golang.org/x/oauth2"
)

func loadTokenFromFile(fname string) (*oauth2.Token, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	tkn := &oauth2.Token{}
	p := json.NewDecoder(f)
	if err := p.Decode(tkn); err != nil {
		return nil, err
	}

	return tkn, nil
}
