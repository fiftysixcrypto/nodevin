package cmd

import (
    "github.com/spf13/viper"
    "log"
)

type Config struct {
    Network      string `mapstructure:"network"`
    StoragePath  string `mapstructure:"storage_path"`
    Port         int    `mapstructure:"port"`
    ResourceLimit string `mapstructure:"resource_limit"`
}

var config Config

func initConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    // Read in config file if it exists
    if err := viper.ReadInConfig(); err == nil {
        log.Println("Using config file:", viper.ConfigFileUsed())
    }

    // Bind CLI flags to viper
    viper.BindPFlag("network", rootCmd.Flags().Lookup("network"))
    viper.BindPFlag("storage_path", rootCmd.Flags().Lookup("storage_path"))
    viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
    viper.BindPFlag("resource_limit", rootCmd.Flags().Lookup("resource_limit"))

    // Unmarshal the config into the config struct
    if err := viper.Unmarshal(&config); err != nil {
        logError("Unable to decode into struct: " + err.Error())
    }
}

