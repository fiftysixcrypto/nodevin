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

func GetLitecoinNetworkComposeConfig(network string) (NetworkConfig, error) {
	baseConfig := NetworkConfig{
		Image:    "fiftysix/litecoin-core",
		Version:  "latest",
		Ports:    []string{"9332:9332", "9333:9333"},
		Volumes:  []string{"litecoin-core-data:/node/litecoin-core"},
		Networks: []string{"litecoin-net"},
		NetworkDefs: map[string]NetworkDetails{
			"litecoin-net": {
				Driver: "bridge",
			},
		},
		VolumeDefs: map[string]VolumeDetails{
			"litecoin-core-data": {
				Labels: map[string]string{
					"nodevin.blockchain.software": "litecoin-core",
				},
			},
		},
	}

	switch network {
	case "litecoin":
		baseConfig.ContainerName = "litecoin-core"
		baseConfig.Command = "litecoind --server=1 --rpcbind=0.0.0.0 --rpcport=9332 --rpcallowip=0.0.0.0/0"
	case "litecoin-testnet":
		baseConfig.ContainerName = "litecoin-core-testnet"
		baseConfig.Command = "litecoind --testnet --server=1 --rpcbind=0.0.0.0 --rpcport=19332 --rpcallowip=0.0.0.0/0"
		baseConfig.Ports = []string{"19332:19332", "19333:19333"}
		baseConfig.Volumes = []string{"litecoin-core-testnet-data:/node/litecoin-core"}
		baseConfig.Networks = []string{"litecoin-testnet-net"}
		baseConfig.NetworkDefs = map[string]NetworkDetails{
			"litecoin-testnet-net": {
				Driver: "bridge",
			},
		}
		baseConfig.VolumeDefs = map[string]VolumeDetails{
			"litecoin-core-testnet-data": {
				Labels: map[string]string{
					"nodevin.blockchain.software": "litecoin-core",
				},
			},
		}
	default:
		return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
	}

	cookieAuth := viper.GetBool("cookie-auth")

	if !cookieAuth {
		// Add RPC user/pass to command
		rpcUsername := viper.GetString("rpc-user")
		rpcPassword := viper.GetString("rpc-pass")

		if rpcUsername == "" {
			rpcUsername = "user"
		}

		if rpcPassword == "" {
			rpcPassword = "fiftysix"
		}

		baseConfig.Command = baseConfig.Command + " " + fmt.Sprintf("-rpcuser=%s", rpcUsername) + " " + fmt.Sprintf("-rpcpassword=%s", rpcPassword)
	}

	return baseConfig, nil
}
