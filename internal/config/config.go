/*
// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2024 The Nodevin Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
*/

package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Add project-specific configuration (current directory)
	viper.AddConfigPath(".")

	// Add user-specific configuration (~/.nodevin)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(filepath.Join(homeDir, ".nodevin"))
	}

	// Add global configuration (/etc/nodevin)
	viper.AddConfigPath("/etc/nodevin")

	// Add executable-specific configuration
	exePath, err := os.Executable()
	if err == nil {
		viper.AddConfigPath(filepath.Dir(exePath))
	}

	// Allow overriding with environment variables
	viper.AutomaticEnv()

	// Read the configuration file, if available
	err = viper.ReadInConfig()
	if err == nil {
		logger.LogInfo(fmt.Sprintf("Successfully loaded configuration from %s", viper.ConfigFileUsed()))
	} else {
		logger.LogError(fmt.Sprintf("Warning: Could not read .env file: %v", err))
	}
}
