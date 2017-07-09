package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	credentials = flag.String("credentials", defaultCredentials, "where to load secrets from")
	portFlag    = flag.String("port", "3030", "port to bind")
	tokenFile   = flag.String("token", "", "load token from this file")
)

func init() {
	flag.Parse()
}

func main() {
	cfg, err := loadConfigFromJSON(*credentials)
	if err != nil {
		if tokenFile == nil || *tokenFile == "" {
			log.Fatal(err)
		} else {
			log.Printf("No secrets available, will attempt token load")
		}
	}

	port := os.Getenv("AI_LIFE_PORT")
	if port == "" {
		port = *portFlag
	}

	lcfg := zap.NewDevelopmentConfig()
	lcfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	lcfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	logger, err := lcfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	appCfg := appConfig{
		tokenFile: tokenFile,
		log:       logger,
	}

	handler := newHandler(cfg, appCfg)

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
	registerSignals(handler)

	logger.Info("Ready to serve ", zap.String("port", port))
	logger.Fatal("Shutting down", zap.Error(http.ListenAndServe(fmt.Sprintf(":%s", port), nil)))
}
