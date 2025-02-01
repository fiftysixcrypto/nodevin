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

package nodes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ContainerInfo struct {
	ID         string `json:"ID"`
	Image      string `json:"Image"`
	Command    string `json:"Command"`
	CreatedAt  string `json:"CreatedAt"`
	RunningFor string `json:"RunningFor"`
	Status     string `json:"Status"`
	Ports      string `json:"Ports"`
	Names      string `json:"Names"`
}

type RPCResponse struct {
	Result interface{} `json:"result"`
	Error  *RPCError   `json:"error"`
	ID     string      `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var infoCmd = &cobra.Command{
	Use:   "info [network]",
	Short: "Get information about running blockchain nodes",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var networkFilter string
		if len(args) > 0 {
			networkFilter = args[0]
		}
		displayInfo(networkFilter)
	},
}

// Prepare the docker ps command
func displayInfo(networkFilter string) {
	args := []string{"ps", "--format", "{{json .}}"}

	// Execute the docker ps command
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogError("Failed to fetch Docker container information: " + err.Error())
		return
	}

	fmt.Println("\n-- Running Nodes:\n")

	// Parse the output
	containers := strings.Split(string(output), "\n")
	if len(containers) < 2 {
		fmt.Println("No running blockchain nodes found.\n")
		displayNodeDirectoryInfo(networkFilter)
		fmt.Println("\n-- Helpful Commands:\n")
		fmt.Printf("%s start <network>\n", utils.GetNodevinExecutable())
		fmt.Printf("%s start <network> --testnet\n", utils.GetNodevinExecutable())
		fmt.Printf("%s stop <network>\n", utils.GetNodevinExecutable())
		return
	}

	// Set up tabwriter for nicely formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "| BLOCKCHAIN\t VERSION\t COMMAND\t STATUS\t PORTS\t PEERS\t LATEST BLOCK")

	for _, containerJSON := range containers {
		if strings.TrimSpace(containerJSON) == "" {
			continue
		}
		var container ContainerInfo
		if err := json.Unmarshal([]byte(containerJSON), &container); err != nil {
			logger.LogError("Failed to parse container JSON: " + err.Error())
			continue
		}

		imageName := container.Image
		if strings.HasPrefix(container.Image, "fiftysix/") {
			imageName = strings.TrimPrefix(container.Image, "fiftysix/")
		} else {
			// Skip if the image does not begin with fiftysix
			continue
		}

		if networkFilter != "" && !strings.Contains(container.Names, networkFilter) {
			continue
		}

		version := "unknown"
		if parts := strings.Split(imageName, ":"); len(parts) > 1 {
			imageName = parts[0]
			version = parts[1]
			if version == "latest" {
				// Fetch the actual version from the container's environment
				version = getNodeVersionFromEnv(container.ID)
			}
		}

		formattedPorts := formatPorts(container.Ports)

		if !utils.IsSupportedExtendedInfoSoftware(imageName) {
			fmt.Fprintf(w, "| %s\t %s\t %s\t %s\t %s\t %s\t %s/%s\n",
				container.Names,
				version,
				container.Command,
				container.Status,
				formattedPorts,
				"-",
				"-",
				"-",
			)
			continue
		}

		localLatestBlock, globalLatestBlock := getLatestBlocks(container.Names)
		peers := getPeers(container.Names)

		fmt.Fprintf(w, "| %s\t %s\t %s\t %s\t %s\t %d\t %d/%d\n",
			container.Names,
			version,
			container.Command,
			container.Status,
			formattedPorts,
			peers,
			localLatestBlock,
			globalLatestBlock,
		)
	}
	w.Flush()

	fmt.Println("")

	displayNodeDirectoryInfo(networkFilter)

	fmt.Println("\n-- Helpful Commands:\n")

	fmt.Printf("%s stop <network>\n", utils.GetNodevinExecutable())
	fmt.Printf("%s shell <network>\n", utils.GetNodevinExecutable())
	fmt.Printf("%s logs <network> --tail 20\n", utils.GetNodevinExecutable())
}

func getLatestBlocks(containerName string) (int, int) {
	localLatestBlock := getLocalLatestBlock(containerName)
	globalLatestBlock := getGlobalLatestBlock(containerName)

	return localLatestBlock, globalLatestBlock
}

func getLocalLatestBlock(containerName string) int {
	url := getLocalEndpointByContainerName(containerName)
	method := "getblockcount"
	params := "[]"
	user := viper.GetString("rpc-user")
	pass := viper.GetString("rpc-pass")

	response, err := makeRequest("", url, method, params, "", user, pass)
	if err != nil {
		//logger.LogError("Failed to get local latest block: " + err.Error())
		return 0
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(response, &rpcResponse); err != nil {
		//logger.LogError("Failed to parse RPC response: " + err.Error())
		return 0
	}

	if rpcResponse.Error != nil {
		//logger.LogError("RPC Error: " + rpcResponse.Error.Message)
		return 0
	}

	blockCount, ok := rpcResponse.Result.(float64)
	if !ok {
		//logger.LogError("Failed to parse block count")
		return 0
	}

	return int(blockCount)
}

func getGlobalLatestBlock(containerName string) int {
	globalFetchLink := getGlobalEndpointByContainerName(containerName)

	resp, err := http.Get(globalFetchLink)
	if err != nil {
		logger.LogError("Failed to fetch global latest block: " + err.Error())
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.LogError("Failed to read response body: " + err.Error())
		return 0
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logger.LogError("Failed to parse response body: " + err.Error())
		return 0
	}

	blockCount, ok := result["height"].(float64)
	if !ok {
		logger.LogError("Failed to parse global block count")
		return 0
	}

	return int(blockCount)
}

func getGlobalEndpointByContainerName(containerName string) string {
	globalFetchLink := ""

	if containerName == "bitcoin-core" {
		globalFetchLink = "https://blockchain.info/latestblock"
	} else if containerName == "bitcoin-core-testnet" {
		globalFetchLink = "https://api.blockcypher.com/v1/btc/test3"
	} else if containerName == "litecoin-core" {
		globalFetchLink = "https://api.blockcypher.com/v1/ltc/main"
	} else if containerName == "litecoin-core-testnet" {
		globalFetchLink = ""
	} else if containerName == "dogecoin-core" {
		globalFetchLink = "https://api.blockcypher.com/v1/doge/main"
	} else if containerName == "dogecoin-core-testnet" {
		globalFetchLink = ""
	}

	return globalFetchLink
}

func getLocalEndpointByContainerName(containerName string) string {
	url := "http://127.0.0.1"

	if containerName == "bitcoin-core" {
		url = "http://127.0.0.1:8332"
	} else if containerName == "bitcoin-core-testnet" {
		url = "http://127.0.0.1:18332"
	} else if containerName == "litecoin-core" {
		url = "http://127.0.0.1:9332"
	} else if containerName == "litecoin-core-testnet" {
		url = "http://127.0.0.1:19332"
	} else if containerName == "dogecoin-core" {
		url = "http://127.0.0.1:22555"
	} else if containerName == "dogecoin-core-testnet" {
		url = "http://127.0.0.1:44555"
	}

	return url
}

func getPeers(containerName string) int {
	url := getLocalEndpointByContainerName(containerName)
	method := "getconnectioncount"
	params := "[]"
	user := viper.GetString("rpc-user")
	pass := viper.GetString("rpc-pass")

	response, err := makeRequest("bitcoin", url, method, params, "", user, pass)
	if err != nil {
		//logger.LogError("Failed to get peer count: " + err.Error())
		return 0
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(response, &rpcResponse); err != nil {
		//logger.LogError("Failed to parse RPC response: " + err.Error())
		return 0
	}

	if rpcResponse.Error != nil {
		//logger.LogError("RPC Error: " + rpcResponse.Error.Message)
		return 0
	}

	peerCount, ok := rpcResponse.Result.(float64)
	if !ok {
		//logger.LogError("Failed to parse peer count")
		return 0
	}

	return int(peerCount)
}

func displayNodeDirectoryInfo(networkFilter string) {
	fmt.Println("-- Blockchain Node Data:\n")

	// Get the nodevin data directory
	nodevinDataDir, err := utils.GetNodevinDataDir()
	if err != nil {
		logger.LogError("Failed to find Nodevin data directory: " + err.Error())
		return
	}

	// Fetch the list of supported networks
	networks := utils.GetAllSupportedNetworks()
	if networks == "" {
		fmt.Println("No data found.")
		return
	}

	// Set up tabwriter for nicely formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "| NETWORK\t SIZE\t DIRECTORY")

	printed := 0

	// Iterate over each supported network and calculate its directory size
	for _, network := range strings.Split(networks, ", ") {
		containerName, exists := utils.GetDefaultLocalMappedContainerName(network)
		if !exists {
			logger.LogError("Unsupported blockchain network: " + network)
			continue
		}

		if networkFilter != "" && networkFilter != network {
			continue
		}

		networkDir := filepath.Join(nodevinDataDir, containerName)
		size, err := getDirectorySize(networkDir)
		sizeDescription := "unknown"
		if err == nil {
			printed++
			sizeDescription = utils.GetSizeDescription(size)
			// Output the formatted row with network name, size, and directory path
			fmt.Fprintf(w, "| %s\t %s\t %s\n", network, sizeDescription, networkDir)
		}
	}

	if printed == 0 {
		fmt.Fprintf(w, "| -\t -\t -\n")
	}

	w.Flush()
}

func getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

func formatPorts(ports string) string {
	portSegments := strings.Split(ports, ",")
	formattedPorts := []string{}
	uniquePorts := make(map[string]bool)

	for _, segment := range portSegments {
		// Use a regex to extract the port numbers and ranges
		re := regexp.MustCompile(`(\d+(-\d+)?)`)
		matches := re.FindAllString(segment, -1)
		for _, match := range matches {
			if match != "0" && !uniquePorts[match] {
				uniquePorts[match] = true
				formattedPorts = append(formattedPorts, match)
			}
		}
	}
	return strings.Join(formattedPorts, ", ")
}

func getNodeVersionFromEnv(containerID string) string {
	// Prepare the docker inspect command
	cmd := exec.Command("docker", "inspect", containerID)

	// Execute the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogError("Failed to inspect container: " + err.Error())
		return "unknown"
	}

	// Parse the output JSON
	var inspectData []map[string]interface{}
	if err := json.Unmarshal(output, &inspectData); err != nil {
		logger.LogError("Failed to parse inspect data: " + err.Error())
		return "unknown"
	}

	// Traverse to find environment variables
	if len(inspectData) > 0 {
		if config, ok := inspectData[0]["Config"].(map[string]interface{}); ok {
			if envVars, ok := config["Env"].([]interface{}); ok {
				for _, envVar := range envVars {
					if envStr, ok := envVar.(string); ok && strings.HasPrefix(envStr, "NODE_VERSION=") {
						return strings.TrimPrefix(envStr, "NODE_VERSION=")
					}
				}
			}
		}
	}

	return "unknown"
}
