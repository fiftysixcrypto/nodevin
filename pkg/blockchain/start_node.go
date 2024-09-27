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

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/fiftysixcrypto/nodevin/pkg/blockchain/bitcoin"
	"github.com/fiftysixcrypto/nodevin/pkg/blockchain/litecoin"
	"github.com/fiftysixcrypto/nodevin/pkg/docker"
	"github.com/fiftysixcrypto/nodevin/pkg/docker/compose"

	"github.com/spf13/cobra"
)

var ord bool

var startNodeCmd = &cobra.Command{
	Use:   "start [network]",
	Short: "Start a blockchain node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startNode(args)
	},
}

func startNode(args []string) {
	if len(args) == 0 {
		logger.LogError("No network provided. Nodevin supports any of the following: " + utils.GetAllSupportedNetworks())
		logger.LogInfo("Example usage: `nodevin start <network>`")
		return
	}

	network := args[0]

	containerName, exists := utils.GetFiftysixDockerhubContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	logger.LogInfo("Starting blockchain node for network: " + network)

	// Initialize Docker client
	if err := docker.InitDockerClient(); err != nil {
		logger.LogError("Failed to initialize Docker client: " + err.Error())
		return
	}

	// Pull the Docker image
	image := containerName + ":latest"
	if err := docker.PullImage(image); err != nil {
		logger.LogError("Failed to pull Docker image: " + err.Error())
		return
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logger.LogError("failed to get current working directory: " + err.Error())
		return
	}

	// Create env file for chain compose
	composeFilePath, err := createComposeFileForNetwork(network, cwd)
	if err != nil {
		logger.LogError("Failed to create node docker compose file: " + err.Error())
	}

	// Start the node
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to start Docker Compose services: " + err.Error())
		return
	}

	logger.LogInfo("Cleaning up excess containers and volumes...")
	if err := compose.RemoveInitContainersAndVolumes(); err != nil {
		logger.LogError("Failed to clean up excess init containers and volumes: " + err.Error())
	} else {
		logger.LogInfo("Successfully cleaned up excess containers and volumes.")
	}

	logger.LogInfo("Successfully started blockchain node for network: " + network)

	startMessage, _ := utils.GetStartMessage(network)

	fmt.Printf("\n%s\n", startMessage)
}

func createComposeFileForNetwork(network string, cwd string) (string, error) {
	switch network {
	case "daemon":
		return "", nil // start the daemon
	case "bitcoin":
		return bitcoin.CreateBitcoinComposeFile(cwd)
	case "ethereum":
		return "", nil // createEthereumComposeFile(cwd)
	case "dogecoin":
		return "", nil // createDogecoinComposeFile(cwd)
	case "ethereumclassic":
		return "", nil // createEthereumClassicComposeFile(cwd)
	case "litecoin":
		return litecoin.CreateLitecoinComposeFile(cwd)
	default:
		return "", fmt.Errorf("unsupported network: %s", network)
	}
}
