package config

import (
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Port      string `mapstructure:"port"`
	DataDir   string `mapstructure:"data-dir"`
	ExtraArgs string `mapstructure:"extra-args"`
	//ResourceLimit string `mapstructure:"resource_limit"`
}

var AppConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Read in config file if it exists
	//if err := viper.ReadInConfig(); err == nil {
	//	log.Println("Using config file:", viper.ConfigFileUsed())
	//}

	// Unmarshal the config into the config struct
	if err := viper.Unmarshal(&AppConfig); err != nil {
		logger.LogError("Unable to decode into struct: " + err.Error())
	}
}
