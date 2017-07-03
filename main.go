package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg, err := loadConfigFromJSON(secretsJSON)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("AI_LIFE_PORT")
	if port == "" {
		port = "3030"
	}

	handler := newHandler(cfg)
	registerServeMux(handler)

	fmt.Println("Ready to serve on port: ", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
