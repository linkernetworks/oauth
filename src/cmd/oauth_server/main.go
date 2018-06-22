package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linkernetworks/logger"
	"github.com/linkernetworks/oauth/src/config"
	"github.com/linkernetworks/oauth/src/server"
)

func main() {

	server := server.New(config.DevelopConfig)
	server.Start()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	s := <-stopSignal
	logger.Infof("Stopped by [%v] signal", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Shutdown OAuth server failed. err: %v", err)
	}
}
