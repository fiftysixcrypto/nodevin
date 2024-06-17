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
	"strings"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/fiftysixcrypto/nodevin/pkg/docker"
	"github.com/spf13/cobra"
)

var deleteVolumeCmd = &cobra.Command{
	Use:   "delete [volume-name-or-image-name]",
	Short: "Delete a Docker volume and its associated images, or a delete specific image",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.LogError("No volume or image name provided. To delete a volume or image, specify the name explicitly.")
			dockerVolumes, err := docker.ListVolumes()
			if err != nil {
				logger.LogError("Unable to list local Docker volumes: " + err.Error())
				return
			}
			logger.LogInfo("List of current volumes: " + dockerVolumes)
			logger.LogInfo("Example usage: `nodevin delete <volume-name-or-image-name>`")
			logger.LogInfo("Example usage: `nodevin delete fiftysix/<image-name>:<tag>`")
			logger.LogInfo("Example usage: `nodevin delete all`")
			return
		}

		name := args[0]

		err := docker.InitDockerClient()
		if err != nil {
			logger.LogError("Failed to initialize Docker client: " + err.Error())
			return
		}

		if name == "all" {
			deleteAllVolumes()
		} else if strings.Contains(name, ":") {
			removeSpecificImage(name)
		} else {
			deleteVolume(name)
		}
	},
}

func deleteVolume(volumeName string) {
	volumes, err := docker.ListVolumes()
	if err != nil {
		logger.LogError("Unable to list local Docker volumes: " + err.Error())
		return
	}

	volumeList := strings.Split(volumes, "\n")
	found := false
	for _, volume := range volumeList {
		if volume == volumeName {
			found = true
			break
		}
	}

	if !found {
		logger.LogError("Volume name not found: " + volumeName)
		logger.LogInfo("List of current volumes: " + volumes)
		logger.LogInfo("Example usage: `nodevin delete <volume-name>`")
		logger.LogInfo("Example usage: `nodevin delete fiftysix/<image-name>:<tag>`")
		logger.LogInfo("Example usage: `nodevin delete all`")
		return
	}

	completed, err := docker.RemoveVolume(volumeName)
	if err != nil {
		logger.LogError("Failed to remove volume: " + err.Error())
		return
	}

	if completed {
		logger.LogInfo("Successfully removed volume: " + volumeName)
		//removeNetworkImages(volumeName)
	} else {
		logger.LogError("Failed to remove volume: " + volumeName)
	}
}

func deleteAllVolumes() {
	logger.LogInfo("Deleting all Docker volumes with label 'nodevin.blockchain.software'...")

	// Run the Docker command to remove all volumes with the nodevin volume label
	removeCmd := exec.Command("sh", "-c", "docker volume ls -q -f label=nodevin.blockchain.software | xargs -r docker volume rm")
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker volumes with label 'nodevin.blockchain.software': " + err.Error())
		return
	}

	logger.LogInfo("Volumes with label 'nodevin.blockchain.software' removed.")

	removeAllImages()
	logger.LogInfo("Successfully deleted all docker volumes.")
}

func removeSpecificImage(image string) {
	logger.LogInfo("Removing specific Docker image: " + image)

	removeCmd := exec.Command("docker", "rmi", image)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker image: " + err.Error())
		return
	}

	logger.LogInfo("Successfully removed Docker image: " + image)
}

func removeNetworkImages(network string) {
	imagePrefix, exists := utils.GetFiftysixDockerhubContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	logger.LogInfo("Removing Docker images for network: " + network)

	removeCmd := exec.Command("sh", "-c", "docker images --format '{{.Repository}}:{{.Tag}}' | grep '^"+imagePrefix+":' | xargs -r docker rmi")
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker images for network " + network + ": " + err.Error())
		return
	}

	logger.LogInfo("Successfully removed Docker images for network: " + network)
}

func removeAllImages() {
	logger.LogInfo("Removing all Docker images starting with 'fiftysix/'...")

	removeCmd := exec.Command("sh", "-c", "docker images --format '{{.Repository}}:{{.Tag}}' | grep '^fiftysix/' | xargs -r docker rmi")
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker images: " + err.Error())
		return
	}

	logger.LogInfo("Successfully removed all Docker images starting with 'fiftysix/'.")
}
