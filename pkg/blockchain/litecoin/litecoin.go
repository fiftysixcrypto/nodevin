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

package litecoin

import (
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/fiftysixcrypto/nodevin/pkg/docker/compose"
	"github.com/spf13/viper"
)

func CreateLitecoinComposeFile(cwd string) (string, error) {
	var network string

	if utils.CheckIfTestnetOrTestnetNetworkFlag() {
		network = "litecoin-testnet"
	} else {
		network = "litecoin"
	}

	litecoinBaseComposeConfig, err := compose.GetLitecoinNetworkComposeConfig(network)
	if err != nil {
		return "", err
	}

	if viper.GetBool("ord") || viper.GetBool("ord-litecoin") {
		var ordLitecoinNetwork string

		if utils.CheckIfTestnetOrTestnetNetworkFlag() {
			ordLitecoinNetwork = "ord-litecoin-testnet"
		} else {
			ordLitecoinNetwork = "ord-litecoin"
		}

		ordLitecoinComposeConfig, err := compose.GetOrdLitecoinNetworkComposeConfig(ordLitecoinNetwork)
		if err != nil {
			return "", err
		}

		composeFilePath, err := compose.CreateComposeFile(
			litecoinBaseComposeConfig.ContainerName,
			litecoinBaseComposeConfig,
			[]string{"ord-litecoin"},
			[]compose.NetworkConfig{ordLitecoinComposeConfig},
			cwd)

		if err != nil {
			return "", err
		}

		return composeFilePath, nil
	} else {
		composeFilePath, err := compose.CreateComposeFile(
			litecoinBaseComposeConfig.ContainerName,
			litecoinBaseComposeConfig,
			[]string{},
			[]compose.NetworkConfig{},
			cwd)

		if err != nil {
			return "", err
		}

		return composeFilePath, nil
	}
}
