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
	"fmt"
	"os"
	"path/filepath"
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
