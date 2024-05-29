package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopNodeCmd = &cobra.Command{
	Use:   "stop-node",
	Short: "Stop a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		stopNode()
	},
}

func stopNode() {
	logInfo("Stopping blockchain node...")

	// Log configuration for debugging
	logInfo(fmt.Sprintf("Network: %s, StoragePath: %s, Port: %d, ResourceLimit: %s", config.Network, config.StoragePath, config.Port, config.ResourceLimit))

	// Get the blockchain instance
	blockchain, exists := GetBlockchain(config.Network)
	if !exists {
		logError("Unsupported blockchain network: " + config.Network)
		return
	}

	// Stop the blockchain node
	if err := blockchain.StopNode(); err != nil {
		logError("Failed to stop blockchain node: " + err.Error())
		return
	}

	logInfo("Blockchain node stopped successfully")
}

func init() {
	stopNodeCmd.Flags().StringVar(&config.Network, "network", "mainnet", "Blockchain network to connect to")

	viper.BindPFlag("network", stopNodeCmd.Flags().Lookup("network"))

	rootCmd.AddCommand(stopNodeCmd)
}
