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
	"path/filepath"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [network-name-or-all]",
	Short: "Delete a directory associated with a network or delete all data",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the nodevin data directory
		nodevinDataDir, err := utils.GetNodevinDataDir()
		if err != nil {
			logger.LogError("Failed to find Nodevin data directory: " + err.Error())
			return
		}

		if len(args) == 0 {
			logger.LogError("No network name provided. To delete network data, specify the name explicitly (for example: bitcoin, litecoin)")
			logger.LogInfo("Example usage: `nodevin delete <network>`")
			logger.LogInfo("Example usage: `nodevin delete all`")
			return
		}

		name := args[0]

		if name == "all" {
			deleteAllDirectories(nodevinDataDir)
		} else {
			deleteNetworkDirectory(nodevinDataDir, name)
		}
	},
}

func deleteNetworkDirectory(baseDir, networkName string) {
	containerName, exists := utils.GetDefaultLocalMappedContainerName(networkName)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + networkName)
		return
	}

	networkDir := filepath.Join(baseDir, containerName)
	if _, err := os.Stat(networkDir); os.IsNotExist(err) {
		logger.LogError("Data for network not found: " + networkDir)
		return
	}

	// Remove the network directory
	err := os.RemoveAll(networkDir)
	if err != nil {
		logger.LogError("Failed to remove data for network " + networkName + ": " + err.Error())
		return
	}

	logger.LogInfo(fmt.Sprintf("Successfully removed %s data directory", networkName))
}

func deleteAllDirectories(baseDir string) {
	// Remove the entire nodevinDataDir directory
	err := os.RemoveAll(baseDir)
	if err != nil {
		logger.LogError("Failed to remove all directories: " + err.Error())
		return
	}

	logger.LogInfo("Successfully removed all nodevin blockchain data")
}