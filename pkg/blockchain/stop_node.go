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

package blockchain

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"

	"github.com/spf13/cobra"
)

var stopNodeCmd = &cobra.Command{
	Use:   "stop [network]",
	Short: "Stop a blockchain node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.LogError("No network specified. To stop a node, specify the network explicitly.")
			logger.LogInfo("")
			availableNetworks := utils.GetAllSupportedNetworks()

			logger.LogInfo("List of available networks: " + availableNetworks)
			logger.LogInfo("Example usage: `nodevin stop <network>`")
			logger.LogInfo("Example usage: `nodevin stop <network> --testnet`")
			logger.LogInfo("Example usage: `nodevin stop <network> --network=\"goerli\"`")
			logger.LogInfo("Example usage: `nodevin stop all`")
			return
		}

		network := args[0]
		if network == "all" {
			stopAllNodes()
		} else {
			stopNode(network)
		}
	},
}

func stopNode(network string) {
	logger.LogInfo("Stopping blockchain node...")

	containerName, exists := utils.GetFiftysixLocalMappedContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	// Get nodevin absolute path
	composeCreatePath, err := os.Executable()
	if err != nil {
		cwd, wdErr := os.Getwd()
		if wdErr != nil {
			composeCreatePath = ""
		} else {
			composeCreatePath = cwd
		}
	}

	// Get the directory where the executable is located
	composeCreateDir := filepath.Dir(composeCreatePath)

	composeFileName := fmt.Sprintf("docker-compose_%s.yml", containerName)
	if utils.CheckIfTestnetOrTestnetNetworkFlag() {
		composeFileName = fmt.Sprintf("docker-compose_%s.yml", containerName+"-testnet")
	}

	composeFilePath := filepath.Join(composeCreateDir, composeFileName)

	// Check if there are any running containers for this compose file
	psCmd := exec.Command("docker-compose", "-f", composeFilePath, "ps", "-q")
	psOut, err := psCmd.Output()
	if err != nil {
		logger.LogError("Failed to find Docker Compose services: " + err.Error())
		return
	}

	if len(psOut) == 0 {
		logger.LogInfo("No running containers found for the specified network (did you mean to add --testnet?)")
		return
	}

	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker Compose services: " + err.Error())
		return
	}

	logger.LogInfo("Blockchain node stopped successfully.")
}

func stopAllNodes() {
	logger.LogInfo("Stopping all Docker containers...")

	psCmd := exec.Command("sh", "-c", "docker ps -q")
	psOut, err := psCmd.Output()
	if err != nil {
		logger.LogError("Failed to list running Docker containers: " + err.Error())
		return
	}

	if len(psOut) == 0 {
		logger.LogInfo("No running Docker containers found.")
		return
	}

	stopCmd := exec.Command("sh", "-c", "docker stop $(docker ps -aq)")
	stopCmd.Stdout = os.Stdout
	stopCmd.Stderr = os.Stderr

	if err := stopCmd.Run(); err != nil {
		logger.LogError("Failed to stop all Docker containers: " + err.Error())
		return
	}

	rmCmd := exec.Command("sh", "-c", "docker rm $(docker ps -aq)")
	rmCmd.Stdout = os.Stdout
	rmCmd.Stderr = os.Stderr

	if err := rmCmd.Run(); err != nil {
		logger.LogError("Failed to remove all Docker containers: " + err.Error())
		return
	}

	logger.LogInfo("All Docker containers stopped and removed successfully.")
}
