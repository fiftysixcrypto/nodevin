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
	"fmt"
	"os"
	"os/exec"

	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/spf13/cobra"
)

var ipfsSupportCmd = &cobra.Command{
	Use:   "ipfs support [network]",
	Short: "Support and pin a network snapshot in a local IPFS container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		CID, exists := utils.GetSnapshotCIDByNetwork(name)
		if !exists {
			fmt.Printf("Unsupported name: %s\n", name)
			return
		}

		containerName := "ipfs"
		err := pinCIDInContainer(containerName, CID)
		if err != nil {
			fmt.Printf("Failed to pin CID in container %s: %v\n", containerName, err)
		} else {
			fmt.Printf("Successfully pinned CID %s for %s in container %s\n", CID, name, containerName)
		}
	},
}

func pinCIDInContainer(containerName, CID string) error {
	args := []string{"exec", "-it", containerName, "ipfs", "pin", "add", CID}

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
