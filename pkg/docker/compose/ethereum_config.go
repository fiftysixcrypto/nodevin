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
)

func GetEthereumNetworkComposeConfig(network string) (NetworkConfig, error) {
	if network == "ethereum" {
		return NetworkConfig{
			Image:         "ethereum/client-go:latest",
			ContainerName: "ethereum-node",
			Command:       "geth",
			Ports:         []string{"8545:8545", "30303:30303"},
			Volumes:       []string{"ethereum-node-data:/root/.ethereum"},
			Networks:      []string{"ethereum-net"},
			NetworkDefs: map[string]NetworkDetails{
				"ethereum-net": {
					Driver: "bridge",
				},
			},
			VolumeDefs: map[string]VolumeDetails{
				"ethereum-node-data": {
					Labels: map[string]string{
						"nodevin.blockchain.software": "ethereum-node",
					},
				},
			},
		}, nil
	}
	return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
}
