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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
)

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

// RemoveInitContainersAndVolumes removes all containers that match "init-config-*" and their associated volumes
func RemoveInitContainersAndVolumes() error {
	// Find all containers that match the pattern "init-config-*"
	logger.LogInfo("Searching for init-config-* containers...")

	listContainersCmd := exec.Command("docker", "ps", "-a", "--filter", "name=init-config-", "--format", "{{.Names}}")
	var containerListOutput bytes.Buffer
	listContainersCmd.Stdout = &containerListOutput
	if err := listContainersCmd.Run(); err != nil {
		logger.LogError("Failed to list containers with name pattern init-config-*: " + err.Error())
		return err
	}

	// Parse the list of container names
	containerNames := strings.Split(strings.TrimSpace(containerListOutput.String()), "\n")
	if len(containerNames) == 0 || containerNames[0] == "" {
		logger.LogInfo("No containers with name pattern init-config-* found.")
		return nil
	}

	// Loop through each container and remove it along with its associated volume
	for _, containerName := range containerNames {
		if containerName == "" {
			continue
		}

		// Stop and remove the container
		logger.LogInfo(fmt.Sprintf("Stopping and removing container: %s", containerName))
		removeContainerCmd := exec.Command("docker", "rm", "-f", containerName)
		removeContainerOutput, err := removeContainerCmd.CombinedOutput()
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed to remove container: %s, Error: %s", containerName, err.Error()))
			continue // Continue even if one container fails
		}
		logger.LogInfo(fmt.Sprintf("Container removed: %s", string(removeContainerOutput)))

		// Attempt to find and remove associated volume
		// Assuming that the volume has a name pattern similar to the container
		volumeName := fmt.Sprintf("%s-init-volume", containerName)
		logger.LogInfo(fmt.Sprintf("Attempting to remove associated volume: %s", volumeName))
		removeVolumeCmd := exec.Command("docker", "volume", "rm", volumeName)
		removeVolumeOutput, err := removeVolumeCmd.CombinedOutput()
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed to remove volume: %s, Error: %s", volumeName, err.Error()))
			continue
		}
		logger.LogInfo(fmt.Sprintf("Volume removed: %s", string(removeVolumeOutput)))
	}

	return nil
}
