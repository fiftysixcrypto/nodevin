package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"

	"github.com/spf13/cobra"
)

var restartNodeCmd = &cobra.Command{
	Use:   "restart [network]",
	Short: "Restart a blockchain node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.LogError("No network specified. To restart a node, specify the network explicitly.")
			logger.LogInfo("")
			availableNetworks := utils.GetAllSupportedNetworks()

			logger.LogInfo("List of available networks: " + availableNetworks)
			logger.LogInfo("Example usage: `nodevin restart <network>`")
			logger.LogInfo("Example usage: `nodevin restart <network> --testnet`")
			logger.LogInfo("Example usage: `nodevin restart <network> --network=\"goerli\"`")
			return
		}

		network := args[0]
		restartNode(network)
	},
}

func restartNode(network string) {
	logger.LogInfo("Restarting blockchain node...")

	containerName, exists := utils.GetFiftysixLocalMappedContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	// Stop the node
	composeFilePath := fmt.Sprintf("docker-compose_%s.yml", containerName)
	if utils.CheckIfTestnetOrTestnetNetworkFlag() {
		composeFilePath = fmt.Sprintf("docker-compose_%s.yml", containerName+"-testnet")
	}

	stopCmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	stopCmd.Stdout = os.Stdout
	stopCmd.Stderr = os.Stderr

	if err := stopCmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker Compose services: " + err.Error())
		return
	}

	// Start the node
	startCmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr

	if err := startCmd.Run(); err != nil {
		logger.LogError("Failed to start Docker Compose services: " + err.Error())
		return
	}

	logger.LogInfo("Blockchain node restarted successfully.")
}
