package main

import "go.uber.org/zap"

type appConfig struct {
	tokenFile *string
	log       *zap.Logger
}
