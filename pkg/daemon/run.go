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
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/update"
)

func runDaemon() {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogError("Failed to open log file: " + err.Error())
		return
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(multiWriter)

	logger.LogInfo("")
	logger.LogInfo("Starting nodevin daemon...")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	updateTicker := time.NewTicker(24 * time.Hour)
	imageCheckTicker := time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-updateTicker.C:
				update.CommandCheckForUpdatesWorkflow()
			case <-imageCheckTicker.C:
				update.CommandCheckAndUpdateDockerImagesWorkflow()
			}
		}
	}()

	logger.LogInfo("Successfully started daemon.")
	logger.LogInfo("Nodevin monitors running nodes and updates them if newer versions are released.")

	sig := <-sigs
	logger.LogInfo("Received signal: " + sig.String())

	logger.LogInfo("Shutting down nodevin daemon...")
}
