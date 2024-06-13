package blockchain

import (
	"fmt"
	"sync"

	"github.com/fiftysixcrypto/nodevin/internal/config"
	"github.com/spf13/cobra"
)

var (
	RequestCmd      = requestCmd
	BackupCmd       = backupCmd
	RestartNodeCmd  = restartNodeCmd
	ShellCmd        = shellCmd
	StartNodeCmd    = startNodeCmd
	StopNodeCmd     = stopNodeCmd
	DeleteVolumeCmd = deleteVolumeCmd
	LogsCmd         = logsCmd
	InfoCmd         = infoCmd
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
