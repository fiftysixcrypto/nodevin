package main

import (
	"github.com/curveballdaniel/nodevin/internal/config"
	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/root"
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
