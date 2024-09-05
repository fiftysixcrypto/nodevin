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

package compose

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetOrdLitecoinNetworkComposeConfig(network string) (NetworkConfig, error) {
	baseConfig := NetworkConfig{
		Image:    "fiftysix/ord-litecoin",
		Version:  "latest",
		Restart:  "always",
		Ports:    []string{"80:80"},
		Volumes:  []string{"litecoin-core-data:/node/litecoin-core", "ord-litecoin-data:/node/ord-litecoin"},
		Networks: []string{"litecoin-net"},
		NetworkDefs: map[string]NetworkDetails{
			"litecoin-net": {
				Driver: "bridge",
			},
		},
		VolumeDefs: map[string]VolumeDetails{
			"ord-litecoin-data": {
				Labels: map[string]string{
					"nodevin.blockchain.software": "ord-litecoin",
				},
			},
		},
	}

	switch network {
	case "ord-litecoin":
		baseConfig.ContainerName = "ord-litecoin"
		baseConfig.Command = "ord --litecoin-rpc-url http://litecoin-core:9332"
	case "ord-litecoin-testnet":
		baseConfig.ContainerName = "ord-litecoin-testnet"
		baseConfig.Command = "ord --testnet --litecoin-rpc-url http://litecoin-core-testnet:19332"
		baseConfig.Volumes = []string{"litecoin-core-testnet-data:/node/litecoin-core", "ord-litecoin-testnet-data:/node/ord-litecoin"}
		baseConfig.Networks = []string{"litecoin-testnet-net"}
		baseConfig.NetworkDefs = map[string]NetworkDetails{
			"litecoin-testnet-net": {
				Driver: "bridge",
			},
		}
		baseConfig.VolumeDefs = map[string]VolumeDetails{
			"ord-litecoin-testnet-data": {
				Labels: map[string]string{
					"nodevin.blockchain.software": "ord-litecoin",
				},
			},
		}
	default:
		return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
	}

	cookieAuth := viper.GetBool("ord-litecoin-cookie-auth") || viper.GetBool("ord-cookie-auth")

	if !cookieAuth {
		// Add RPC user/pass to command
		rpcUsername := viper.GetString("ord-litecoin-rpc-user")
		rpcPassword := viper.GetString("ord-litecoin-rpc-pass")

		if rpcUsername == "" {
			rpcUsername = viper.GetString("ord-rpc-user")
			if rpcUsername == "" {
				rpcUsername = "user"
			}
		}

		if rpcPassword == "" {
			rpcPassword = viper.GetString("ord-rpc-pass")
			if rpcPassword == "" {
				rpcPassword = "fiftysix"
			}
		}

		baseConfig.Command = baseConfig.Command + " " + fmt.Sprintf("--litecoin-rpc-username %s", rpcUsername) + " " + fmt.Sprintf("--litecoin-rpc-password %s", rpcPassword)
	}

	baseConfig.Command = baseConfig.Command + " server"

	return baseConfig, nil
}
