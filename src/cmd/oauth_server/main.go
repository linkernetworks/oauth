package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/linkernetworks/logger"
)

func init() {
	cf := logger.LoggerConfig{
		Dir:           "./logs",
		Level:         "debug",
		MaxAge:        "720h",
		SuffixPattern: ".%Y%m%d",
		LinkName:      "oauth_log",
	}
	logger.Setup(cf)
}

func main() {

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

	s := <-stopSignal
	logger.Infof("Stopped by [%v] signal", s.String())
}
