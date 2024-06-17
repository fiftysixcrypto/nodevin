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
	"path/filepath"
	"strings"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/docker"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup [volume-name] [destination-path]",
	Short: "Backup a Docker volume",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		volumeName := args[0]
		destinationPath := args[1]

		err := docker.InitDockerClient()
		if err != nil {
			logger.LogError("Failed to initialize Docker client: " + err.Error())
			return
		}

		backupVolume(volumeName, destinationPath)
	},
}

func volumeExists(volumeName string) (bool, error) {
	cmd := exec.Command("docker", "volume", "ls", "-q", "--filter", "name="+volumeName)
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	volumes := strings.Split(string(output), "\n")
	for _, volume := range volumes {
		if volume == volumeName {
			return true, nil
		}
	}
	return false, nil
}

func backupVolume(volumeName string, destinationPath string) {
	logger.LogInfo("Backing up Docker volume...")

	exists, err := volumeExists(volumeName)
	if err != nil {
		logger.LogError("Failed to check if volume exists: " + err.Error())
		return
	}

	if !exists {
		logger.LogError("Specified volume does not exist: " + volumeName)
		return
	}

	absoluteDestinationPath, err := filepath.Abs(destinationPath)
	if err != nil {
		logger.LogError("Failed to resolve absolute destination path: " + err.Error())
		return
	}

	cmd := exec.Command("docker", "run", "--rm", "-v", volumeName+":/volume", "-v", absoluteDestinationPath+":/backup", "alpine", "sh", "-c", "tar czf /backup/backup.tar.gz -C /volume .")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				logger.LogError("Failed to backup volume: " + err.Error() + ". You may need greater permissions to execute this command. Try running with `sudo`.")
			} else {
				logger.LogError("Failed to backup volume: " + err.Error())
			}
		} else {
			logger.LogError("Failed to backup volume: " + err.Error())
		}
		return
	}

	logger.LogInfo("Successfully backed up volume: " + volumeName + " to " + destinationPath)
}
