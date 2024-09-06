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

package update

import (
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var (
	UpdateCmd = updateCmd
)

var updateCmd = &cobra.Command{
	Use:   "update [target]",
	Short: "Update Nodevin software or Docker images",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Update Nodevin software
			CommandCheckForUpdatesWorkflow()
		} else if args[0] == "docker" {
			// Update Docker images
			CommandCheckAndUpdateDockerImagesWorkflow()
		} else {
			logger.LogError("Invalid target specified. Use 'nodevin update' or 'nodevin update docker'.")
		}
	},
}

func CommandCheckForUpdatesWorkflow() {
	logger.LogInfo("Checking for nodevin updates...")

	updateNeeded, err := CheckForUpdates()
	if err != nil {
		logger.LogError("Failed to check for updates: " + err.Error())
		return
	}
	if updateNeeded {
		logger.LogInfo("Update downloaded. Applying update...")
		if err := ApplyUpdate(); err != nil {
			logger.LogError("Failed to apply update: " + err.Error())
			return
		}
		logger.LogInfo("Update applied successfully. Please restart the application.")
	} else {
		logger.LogInfo("nodevin is up to date.")
	}
}

func CommandCheckAndUpdateDockerImagesWorkflow() {
	logger.LogInfo("Checking for Docker image updates...")
	if err := CheckAndUpdateDockerImages(); err != nil {
		logger.LogError("Failed to check/update Docker images: " + err.Error())
		return
	}
	logger.LogInfo("Image updates complete.")
}
