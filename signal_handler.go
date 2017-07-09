package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func registerSignals(h *myHandler) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		h.log.Info("Got signal, saving token", zap.Any("signal", sig))
		h.maybeStoreToken(h.token)
		os.Exit(0)
	}()
}
