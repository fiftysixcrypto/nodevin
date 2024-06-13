package main

import (
	"github.com/fiftysixcrypto/nodevin/internal/config"
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/daemon"
)

func main() {
	// Initialize configurations and loggers
	config.InitConfig()
	logger.Init()

	// Execute the daemon command
	if err := daemon.Execute(); err != nil {
		logger.LogError(err.Error())
	}
}
