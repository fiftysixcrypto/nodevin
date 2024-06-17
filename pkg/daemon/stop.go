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

package daemon

import (
	"os"
	"strconv"
	"syscall"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the nodevin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		stopDaemon()
	},
}

func stopDaemon() {
	pidData, err := os.ReadFile(pidFilePath)
	if err != nil {
		logger.LogError("Failed to read PID file: " + err.Error())
		logger.LogInfo("Are you sure the daemon is running?")
		return
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		logger.LogError("Invalid PID: " + err.Error())
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		logger.LogError("Failed to find process: " + err.Error())
		return
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		logger.LogError("Failed to stop daemon: " + err.Error())
		return
	}

	err = os.Remove(pidFilePath)
	if err != nil {
		logger.LogError("Failed to remove PID file: " + err.Error())
		return
	}

	logger.LogInfo("Nodevin daemon successfully stopped.")
}
