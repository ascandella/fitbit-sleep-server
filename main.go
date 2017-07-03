package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
)

var (
	credentials = flag.String("credentials", defaultCredentials, "where to load secrets from")
	portFlag    = flag.String("port", "3030", "port to bind")
	tokenFile   = flag.String("token", "", "load token from this file for testing")
)

func init() {
	flag.Parse()
}

func main() {
	cfg, err := loadConfigFromJSON(*credentials)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("AI_LIFE_PORT")
	if port == "" {
		port = *portFlag
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	handler := newHandler(cfg, logger)

	if tokenFile != nil && *tokenFile != "" {
		logger.Info("Loading token from file", zap.String("location", *tokenFile))

		token, err := loadTokenFromFile(*tokenFile)
		if err != nil {
			logger.Error("Unable to load token from file", zap.Error(err))
		} else {
			handler.registerToken(token)
		}
	}

	registerServeMux(handler)

	logger.Info("Ready to serve ", zap.String("port", port))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
