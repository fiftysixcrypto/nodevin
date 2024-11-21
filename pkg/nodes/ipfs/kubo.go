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

package ipfs

import (
	"runtime"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/docker"
	"github.com/fiftysixcrypto/nodevin/pkg/docker/compose"
	"github.com/spf13/viper"
)

func CreateKuboComposeFile(cwd string) (string, error) {
	var network = "ipfs"

	kuboBaseComposeConfig, err := compose.GetKuboNetworkComposeConfig(network)
	if err != nil {
		return "", err
	}

	composeFilePath, err := compose.CreateComposeFile(
		kuboBaseComposeConfig.ContainerName,
		kuboBaseComposeConfig,
		[]string{},
		[]compose.NetworkConfig{},
		cwd)

	if err != nil {
		return "", err
	}

	if viper.GetBool("ipfs-cluster") {
		if runtime.GOARCH == "arm64" {
			logger.LogError("Running on ARM architecture: ipfs-cluster functionality is not supported on ARM builds. Skipping ipfs-cluster.")
			return composeFilePath, nil
		}

		ipfsClusterNetwork := "ipfs-cluster"

		// Pull the ipfs-cluster Docker image
		image := viper.GetString("ipfs-cluster-image") + ":" + viper.GetString("ipfs-cluster-version")
		if err := docker.PullImage(image); err != nil {
			logger.LogError("Failed to pull Docker image: " + err.Error())
			return "", err
		}

		ipfsClusterBaseComposeConfig, err := compose.GetIpfsClusterNetworkComposeConfig(ipfsClusterNetwork)
		if err != nil {
			return "", err
		}

		composeFilePath, err := compose.CreateComposeFile(
			kuboBaseComposeConfig.ContainerName,
			kuboBaseComposeConfig,
			[]string{"ipfs-cluster"},
			[]compose.NetworkConfig{ipfsClusterBaseComposeConfig},
			cwd)

		if err != nil {
			return "", err
		}

		return composeFilePath, nil
	} else {
		return composeFilePath, nil
	}
}
