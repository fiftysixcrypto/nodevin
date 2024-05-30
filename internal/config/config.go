package config

import (
	"log"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Network       string `mapstructure:"network"`
	StoragePath   string `mapstructure:"storage_path"`
	Port          int    `mapstructure:"port"`
	ResourceLimit string `mapstructure:"resource_limit"`
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Read in config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Unmarshal the config into the config struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		logger.LogError("Unable to decode into struct: " + err.Error())
	}
}
