package main

import (
	"encoding/json"
	"io/ioutil"

	"golang.org/x/oauth2"
)

func loadTokenFromFile(fname string) (*oauth2.Token, error) {
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	tkn := &oauth2.Token{}
	if err := json.Unmarshal(f, tkn); err != nil {
		return nil, err
	}

	return tkn, nil
}
