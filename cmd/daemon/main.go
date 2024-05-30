package main

import (
	"github.com/curveballdaniel/nodevin/internal/config"
	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/daemon"
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
