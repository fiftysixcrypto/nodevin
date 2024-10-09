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

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/viper"
)

type NetworkInfo struct {
	ContainerName string
	RPCPort       int
	DataSize      int
	SnapshotSize  int
	StartMessage  string
}

var networkInfoMap = map[string]NetworkInfo{
	"bitcoin": {
		ContainerName: "bitcoin-core",
		RPCPort:       8332,
		DataSize:      0,
		SnapshotSize:  0,
		StartMessage:  "\"A system for electronic transactions without relying on trust.\" -- Satoshi Nakamoto",
	},
	"bitcoin-testnet": {
		ContainerName: "bitcoin-core-testnet",
		RPCPort:       18332,
		DataSize:      0,
		SnapshotSize:  0,
		StartMessage:  "\"Testing is the lifeblood of innovation and security.\"",
	},
	"ord": {
		ContainerName: "ord",
		RPCPort:       80,
		DataSize:      0,
		SnapshotSize:  0,
		StartMessage:  "\"Ordinal theory imbues satoshis with numismatic value, allowing them to be collected and traded as curios.\"",
	},
	"ethereum": {
		ContainerName: "geth",
		RPCPort:       8545,
		StartMessage:  "\"Ethereum is the foundation for our digital future.\" -- Vitalik Buterin",
	},
	"ethereum-classic": {
		ContainerName: "core-geth",
		RPCPort:       8545,
		StartMessage:  "\"Code is law.\" -- Ethereum Classic",
	},
	"litecoin": {
		ContainerName: "litecoin-core",
		RPCPort:       9332,
		DataSize:      268435456000, // 250 GB
		SnapshotSize:  429496729600, // 400 GB (~240GB + ~160GB)
		StartMessage:  "\"Litecoin is the silver to Bitcoin's gold.\" -- Charlie Lee",
	},
	"litecoin-testnet": {
		ContainerName: "litecoin-core-testnet",
		RPCPort:       19332,
		DataSize:      0,
		SnapshotSize:  0,
		StartMessage:  "\"Testing is the lifeblood of innovation and security.\"",
	},
	"ord-litecoin": {
		ContainerName: "ord-litecoin",
		RPCPort:       80,
		DataSize:      0,
		SnapshotSize:  0,
		StartMessage:  "\"Ordinal theory imbues satoshis with numismatic value, allowing them to be collected and traded as curios.\"",
	},
	"dogecoin": {
		ContainerName: "dogecoin-core",
		RPCPort:       22555,
		StartMessage:  "\"To the moon!\" -- Dogecoin Community",
	},
}

func NetworkContainerMap() map[string]string {
	networkContainerMap := make(map[string]string)
	for key, value := range networkInfoMap {
		networkContainerMap[key] = value.ContainerName
	}
	return networkContainerMap
}

func NetworkDefaultRPCPorts() map[string]int {
	networkDefaultRPCPorts := make(map[string]int)
	for key, value := range networkInfoMap {
		networkDefaultRPCPorts[key] = value.RPCPort
	}
	return networkDefaultRPCPorts
}

func GetStartMessage(network string) (string, bool) {
	networkInfo, exists := networkInfoMap[network]
	return networkInfo.StartMessage, exists
}

func GetDefaultLocalMappedContainerName(network string) (string, bool) {
	networkInfo, exists := networkInfoMap[network]
	return networkInfo.ContainerName, exists
}

func GetNetworkRequiredDataSize(network string) (int, bool) {
	networkInfo, exists := networkInfoMap[network]
	return networkInfo.DataSize, exists
}

func GetNetworkRequiredSnapshotSize(network string) (int, bool) {
	networkInfo, exists := networkInfoMap[network]
	return networkInfo.SnapshotSize, exists
}

func GetFiftysixDockerhubContainerName(network string) (string, bool) {
	networkInfo, exists := networkInfoMap[network]
	return "fiftysix/" + networkInfo.ContainerName, exists
}

func GetAllSupportedNetworks() string {
	var keys []string
	for key := range networkInfoMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return strings.Join(keys, ", ")
}

func CheckIfTestnetOrTestnetNetworkFlag() bool {
	networkFlag := viper.GetString("network")
	testnetFlag := viper.GetBool("testnet")

	return testnetFlag || networkFlag == "testnet"
}

// Returns path to the user's nodevin data directory (~/.nodevin/data)
func GetNodevinDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	nodevinDataDir := filepath.Join(homeDir, ".nodevin", "data")

	// Create the directory if it doesn't exist
	if _, err := os.Stat(nodevinDataDir); os.IsNotExist(err) {
		err = os.MkdirAll(nodevinDataDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create nodevin data directory: %v", err)
		}
	}

	return nodevinDataDir, nil
}

func GetSizeDescription(size int64) string {
	if size <= 0 {
		return "unknown (do you have proper permissions?)"
	}

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
		PB = TB * 1024
	)

	switch {
	case size >= PB:
		return fmt.Sprintf("%.2f PB", float64(size)/PB)
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
