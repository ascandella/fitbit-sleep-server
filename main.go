package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := loadConfigFromJSON(secretsJSON)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

	state := randStringRunes(24)
	redir := cfg.AuthCodeURL(state)

	handler := myHandler{}
	registerServeMux(myHandler)
	// TODO
	select {}
}
