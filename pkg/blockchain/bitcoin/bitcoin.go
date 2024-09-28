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

package bitcoin

import (
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/fiftysixcrypto/nodevin/pkg/docker"
	"github.com/fiftysixcrypto/nodevin/pkg/docker/compose"
	"github.com/spf13/viper"
)

func CreateBitcoinComposeFile(cwd string) (string, error) {
	var network string

	if utils.CheckIfTestnetOrTestnetNetworkFlag() {
		network = "bitcoin-testnet"
	} else {
		network = "bitcoin"
	}

	bitcoinBaseComposeConfig, err := compose.GetBitcoinNetworkComposeConfig(network)
	if err != nil {
		return "", err
	}

	if viper.GetBool("ord") {
		var ordNetwork string

		if utils.CheckIfTestnetOrTestnetNetworkFlag() {
			ordNetwork = "ord-testnet"
		} else {
			ordNetwork = "ord"
		}

		// Pull the ord Docker image
		image := viper.GetString("ord-image") + ":" + viper.GetString("ord-version")
		if err := docker.PullImage(image); err != nil {
			logger.LogError("Failed to pull Docker image: " + err.Error())
			return "", err
		}

		ordComposeConfig, err := compose.GetOrdNetworkComposeConfig(ordNetwork)
		if err != nil {
			return "", err
		}

		composeFilePath, err := compose.CreateComposeFile(
			bitcoinBaseComposeConfig.ContainerName,
			bitcoinBaseComposeConfig,
			[]string{"ord"},
			[]compose.NetworkConfig{ordComposeConfig},
			cwd)

		if err != nil {
			return "", err
		}

		return composeFilePath, nil
	} else {
		composeFilePath, err := compose.CreateComposeFile(
			bitcoinBaseComposeConfig.ContainerName,
			bitcoinBaseComposeConfig,
			[]string{},
			[]compose.NetworkConfig{},
			cwd)

		if err != nil {
			return "", err
		}

		return composeFilePath, nil
	}
}
