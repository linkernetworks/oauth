package log

import "github.com/linkernetworks/logger"

func Init(config logger.LoggerConfig) error {
	logger.Setup(config)
	return nil
}
