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

package bitcoin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fiftysixcrypto/nodevin/internal/config"
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/blockchain"
)

type Bitcoin struct{}

func (b Bitcoin) StartNode(config config.Config) error {
	fmt.Println("Starting Bitcoin node with config:", config)

	// Define the path to the Docker Compose file
	composeFilePath := "docker/bitcoin/docker-compose_bitcoin-core.yml"

	// Start the services using Docker Compose
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	cmd.Stdout = logWriter{}
	cmd.Stderr = logWriter{}

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to start Docker Compose services: " + err.Error())
		return err
	}

	logger.LogInfo("Bitcoin node started successfully")

	return nil
}

func (b Bitcoin) StopNode() error {
	fmt.Println("Stopping Bitcoin node")

	// Define the path to the Docker Compose file
	composeFilePath := "docker/bitcoin/docker-compose_bitcoin-core.yml"

	// Stop the services using Docker Compose
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	cmd.Stdout = logWriter{}
	cmd.Stderr = logWriter{}

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker Compose services: " + err.Error())
		return err
	}

	logger.LogInfo("Bitcoin node stopped successfully")

	return nil
}

func init() {
	blockchain.RegisterBlockchain("bitcoin", Bitcoin{})
}

type logWriter struct{}

func (f logWriter) Write(bytes []byte) (int, error) {
	return fmt.Fprint(os.Stdout, string(bytes))
}
