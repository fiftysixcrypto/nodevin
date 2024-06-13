package main

import (
	"github.com/fiftysixcrypto/nodevin/internal/config"
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/root"
)

func main() {
	// Initialize loggers
	logger.Init()

	// Initialize configurations
	config.InitConfig()

	// Execute the root command
	if err := root.Execute(); err != nil {
		logger.LogError(err.Error())
	}
}
