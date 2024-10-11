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
	"os/exec"
	"strings"
	"time"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
)

// RemoveInitContainersAndVolumes removes all containers that match "init-config-*" and their associated volumes,
// deletes volumes with the label "nodevin.init.volume", and anonymous volumes created within the last minute.
func RemoveInitContainersAndVolumes() error {
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
		//
	} else {
		// Loop through each container and remove it along with its associated volume
		for _, containerName := range containerNames {
			if containerName == "" {
				continue
			}

			// Stop and remove the container
			removeContainerCmd := exec.Command("docker", "rm", "-f", containerName)
			_, err := removeContainerCmd.CombinedOutput()
			if err != nil {
				logger.LogError(fmt.Sprintf("Failed to remove container: %s, Error: %s", containerName, err.Error()))
				continue // Continue even if one container fails
			}
		}
	}

	// Delete volumes with the label "nodevin.init.volume"
	listInitVolumesCmd := exec.Command("docker", "volume", "ls", "--filter", "label=nodevin.init.volume=true", "--format", "{{.Name}}")
	var initVolumeListOutput bytes.Buffer
	listInitVolumesCmd.Stdout = &initVolumeListOutput
	if err := listInitVolumesCmd.Run(); err != nil {
		logger.LogError("Failed to list volumes with label 'nodevin.init.volume': " + err.Error())
		return err
	}

	// Parse the list of init volumes and delete them
	initVolumeNames := strings.Split(strings.TrimSpace(initVolumeListOutput.String()), "\n")
	for _, volumeName := range initVolumeNames {
		if volumeName == "" {
			continue
		}
		removeVolumeCmd := exec.Command("docker", "volume", "rm", volumeName)
		_, err := removeVolumeCmd.CombinedOutput()
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed to remove init volume: %s, Error: %s", volumeName, err.Error()))
			continue
		}
	}

	// Delete anonymous volumes created within the last minute
	listAnonymousVolumesCmd := exec.Command("docker", "volume", "ls", "--filter", "label=com.docker.volume.anonymous", "--format", "{{.Name}}")
	var anonymousVolumeListOutput bytes.Buffer
	listAnonymousVolumesCmd.Stdout = &anonymousVolumeListOutput
	if err := listAnonymousVolumesCmd.Run(); err != nil {
		logger.LogError("Failed to list anonymous volumes: " + err.Error())
		return err
	}

	anonymousVolumeNames := strings.Split(strings.TrimSpace(anonymousVolumeListOutput.String()), "\n")
	currentTime := time.Now()

	// Loop through each anonymous volume and check the CreatedAt time
	for _, volumeName := range anonymousVolumeNames {
		if volumeName == "" {
			continue
		}

		// Get details of the volume including CreatedAt time
		inspectVolumeCmd := exec.Command("docker", "volume", "inspect", volumeName, "--format", "{{.CreatedAt}}")
		var inspectVolumeOutput bytes.Buffer
		inspectVolumeCmd.Stdout = &inspectVolumeOutput
		if err := inspectVolumeCmd.Run(); err != nil {
			logger.LogError(fmt.Sprintf("Failed to inspect volume: %s, Error: %s", volumeName, err.Error()))
			continue
		}

		// Parse the CreatedAt time
		createdAtStr := strings.TrimSpace(inspectVolumeOutput.String())
		createdAtTime, err := time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed to parse CreatedAt time for volume: %s, Error: %s", volumeName, err.Error()))
			continue
		}

		// Check if the volume was created within the last minute
		if currentTime.Sub(createdAtTime).Minutes() <= 1 {
			removeVolumeCmd := exec.Command("docker", "volume", "rm", volumeName)
			_, err := removeVolumeCmd.CombinedOutput()
			if err != nil {
				logger.LogError(fmt.Sprintf("Failed to remove anonymous volume: %s, Error: %s", volumeName, err.Error()))
				continue
			}
		}
	}

	return nil
}
