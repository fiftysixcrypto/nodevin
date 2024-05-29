package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startNodeCmd = &cobra.Command{
	Use:   "start-node",
	Short: "Start a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		startNode()
	},
}

func startNode() {
	logInfo("Starting blockchain node...")

	// Log configuration for debugging
	logInfo(fmt.Sprintf("Network: %s, StoragePath: %s, Port: %d, ResourceLimit: %s", config.Network, config.StoragePath, config.Port, config.ResourceLimit))

	// Get the blockchain instance
	blockchain, exists := GetBlockchain(config.Network)
	if !exists {
		logError("Unsupported blockchain network: " + config.Network)
		return
	}

	// Start the blockchain node
	if err := blockchain.StartNode(config); err != nil {
		logError("Failed to start blockchain node: " + err.Error())
		return
	}

	logInfo("Blockchain node started successfully")
}

func init() {
	startNodeCmd.Flags().StringVar(&config.Network, "network", "mainnet", "Blockchain network to connect to")
	startNodeCmd.Flags().StringVar(&config.StoragePath, "storage_path", "/var/lib/nodevin", "Path to store blockchain data")
	startNodeCmd.Flags().IntVar(&config.Port, "port", 30303, "Port to bind the node")
	startNodeCmd.Flags().StringVar(&config.ResourceLimit, "resource_limit", "2GB", "Resource limit for the node")

	viper.BindPFlag("network", startNodeCmd.Flags().Lookup("network"))
	viper.BindPFlag("storage_path", startNodeCmd.Flags().Lookup("storage_path"))
	viper.BindPFlag("port", startNodeCmd.Flags().Lookup("port"))
	viper.BindPFlag("resource_limit", startNodeCmd.Flags().Lookup("resource_limit"))

	rootCmd.AddCommand(startNodeCmd)
}
