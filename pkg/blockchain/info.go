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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const pidFilePath = "/tmp/nodevin_daemon.pid"

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
	Use:   "info",
	Short: "Get information about all running blockchain nodes",
	Run: func(cmd *cobra.Command, args []string) {
		displayInfo()
	},
}

func displayInfo() {
	// Prepare the docker ps command
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
		displayVolumeInfo()

		displayPIDInfo()

		fmt.Println("\n-- Helpful Commands:\n")
		fmt.Println("nodevin start <network>")
		fmt.Println("nodevin start <network> --testnet")
		fmt.Println("nodevin delete <volume-name-or-image-name>\n")
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

		version := "unknown"
		if parts := strings.Split(imageName, ":"); len(parts) > 1 {
			imageName = parts[0]
			version = parts[1]
		}

		formattedPorts := formatPorts(container.Ports)

		localLatestBlock := 0
		globalLatestBlock := 0
		peers := 0

		localLatestBlock, globalLatestBlock = getLatestBlocks(container.Names)
		peers = getPeers(container.Names)

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

	displayVolumeInfo()

	displayPIDInfo()

	fmt.Println("\n-- Helpful Commands:\n")

	fmt.Println("nodevin stop <network>")
	fmt.Println("nodevin shell <network>")
	fmt.Println("nodevin daemon start")
	fmt.Println("nodevin logs <network> --tail 50\n")
}

func displayPIDInfo() {
	if _, err := os.Stat(pidFilePath); err == nil {
		fmt.Println("\n-- Nodevin Daemon:")
		pidData, err := os.ReadFile(pidFilePath)
		if err != nil {
			logger.LogError("Failed to read PID file: " + err.Error())
			return
		}
		fmt.Printf("\nDaemon is running to update nodes with PID: %s\n", string(pidData))
		fmt.Println("Logs: nodevin daemon logs")
		fmt.Println("Stop: nodevin daemon stop")
	}
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
		logger.LogError("Failed to get local latest block: " + err.Error())
		return 0
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(response, &rpcResponse); err != nil {
		logger.LogError("Failed to parse RPC response: " + err.Error())
		return 0
	}

	if rpcResponse.Error != nil {
		logger.LogError("RPC Error: " + rpcResponse.Error.Message)
		return 0
	}

	blockCount, ok := rpcResponse.Result.(float64)
	if !ok {
		logger.LogError("Failed to parse block count")
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
		logger.LogError("Failed to get peer count: " + err.Error())
		return 0
	}

	var rpcResponse RPCResponse
	if err := json.Unmarshal(response, &rpcResponse); err != nil {
		logger.LogError("Failed to parse RPC response: " + err.Error())
		return 0
	}

	if rpcResponse.Error != nil {
		logger.LogError("RPC Error: " + rpcResponse.Error.Message)
		return 0
	}

	peerCount, ok := rpcResponse.Result.(float64)
	if !ok {
		logger.LogError("Failed to parse peer count")
		return 0
	}

	return int(peerCount)
}

func displayVolumeInfo() {
	fmt.Println("-- Volume Data:\n")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.LogError("Failed to create Docker client: " + err.Error())
		return
	}

	volumeList, err := cli.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		logger.LogError("Failed to list Docker volumes: " + err.Error())
		return
	}

	if len(volumeList.Volumes) < 1 {
		fmt.Println("No blockchain node data found.")
		return
	}

	// Set up tabwriter for nicely formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "| VOLUME NAME\t SIZE\t MOUNTPOINT")

	for _, vol := range volumeList.Volumes {
		size, err := getVolumeSize(vol.Mountpoint)
		sizeDescription := "unknown"
		if err == nil {
			sizeDescription = getSizeDescription(size)
		} else {
			if err.Error() == "exit status 1" {
				fmt.Println("Failed to get size for volume " + vol.Name + ", do you have proper permissions?")
			} else if err.Error() == "exit status 64" {
				fmt.Println("Failed to get size for volume " + vol.Name + ", do you have proper permissions?")
			} else {
				logger.LogError("Failed to get size for volume " + vol.Name + ": " + err.Error())
			}
		}

		fmt.Fprintf(w, "| %s\t %s\t %s\n",
			vol.Name,
			sizeDescription,
			vol.Mountpoint,
		)
	}
	w.Flush()
}

func getVolumeSize(mountpoint string) (int64, error) {
	cmd := exec.Command("du", "-sb", mountpoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Parse the output to get the size in bytes
	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		return 0, fmt.Errorf("unexpected du output: %s", string(output))
	}

	size, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
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

func getSizeDescription(size int64) string {
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
