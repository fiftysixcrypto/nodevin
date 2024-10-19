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
	"fmt"
	"sync"

	"github.com/fiftysixcrypto/nodevin/internal/config"
	"github.com/spf13/cobra"
)

var (
	RequestCmd   = requestCmd
	ShellCmd     = shellCmd
	StartNodeCmd = startNodeCmd
	StopNodeCmd  = stopNodeCmd
	DeleteCmd    = deleteCmd
	CleanupCmd   = cleanupCmd
	LogsCmd      = logsCmd
	InfoCmd      = infoCmd
	ListCmd      = listCmd
)

var blockchainCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Manage blockchain nodes",
}

func Execute() error {
	return blockchainCmd.Execute()
}

type Blockchain interface {
	StartNode(config.Config) error
	StopNode() error
}

var (
	blockchains   = make(map[string]Blockchain)
	blockchainMtx sync.RWMutex
)

func RegisterBlockchain(name string, blockchain Blockchain) {
	blockchainMtx.Lock()
	defer blockchainMtx.Unlock()

	if _, exists := blockchains[name]; exists {
		panic(fmt.Sprintf("blockchain already registered: %s", name))
	}

	blockchains[name] = blockchain
}

func GetBlockchain(name string) (Blockchain, bool) {
	blockchainMtx.RLock()
	defer blockchainMtx.RUnlock()

	blockchain, exists := blockchains[name]
	return blockchain, exists
}
