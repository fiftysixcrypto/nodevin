package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var stopNodeCmd = &cobra.Command{
	Use:   "stop [network]",
	Short: "Stop a blockchain node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.LogError("No network specified. To stop a node, specify the network explicitly.")
			availableNetworks := getAllSupportedNetworks()

			logger.LogInfo("List of available networks: " + availableNetworks)
			logger.LogInfo("Example usage: `nodevin stop <network>`")
			return
		}

		// TODO: add (no nodes are running!)

		network := args[0]
		stopNode(network)
	},
}

func stopNode(network string) {
	logger.LogInfo("Stopping blockchain node...")

	containerName, exists := getFiftysixLocalMappedContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	// Stop the node
	composeFilePath := fmt.Sprintf("docker/%s/docker-compose_%s.yml", network, containerName)
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker Compose services: " + err.Error())
	}
}
