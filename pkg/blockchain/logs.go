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
	"os"
	"os/exec"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/spf13/cobra"
)

var follow bool
var tail string

var logsCmd = &cobra.Command{
	Use:   "logs [network]",
	Short: "Fetch logs from a running blockchain node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.LogError("No network specified. To fetch logs, specify the network explicitly.")
			availableNetworks := utils.GetAllSupportedNetworks()
			logger.LogInfo("List of available networks: " + availableNetworks)
			logger.LogInfo("Example usage: `nodevin logs <network>`")
			return
		}

		network := args[0]
		fetchLogs(network)
	},
}

func fetchLogs(network string) {
	logger.LogInfo("Fetching logs for blockchain node...")

	properNetwork := network

	if utils.CheckIfTestnetOrTestnetNetworkFlag() {
		if network == "bitcoin" {
			properNetwork = "bitcoin-testnet"
		} else if network == "litecoin" {
			properNetwork = "litecoin-testnet"
		}
	}

	containerName, exists := getOutputLogsContainerName(properNetwork)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	if tail != "" {
		args = append(args, "--tail", tail)
	}
	args = append(args, containerName)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to fetch Docker logs: " + err.Error())
	}
}

func getOutputLogsContainerName(network string) (string, bool) {
	for net, container := range utils.NetworkContainerMap() {
		if net == network {
			return container, true
		}
	}
	return "", false
}

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVar(&tail, "tail", "all", "Number of lines to show from the end of the logs")
}
