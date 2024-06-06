package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/internal/utils"

	"github.com/spf13/cobra"
)

var stopNodeCmd = &cobra.Command{
	Use:   "stop [network]",
	Short: "Stop a blockchain node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.LogError("No network specified. To stop a node, specify the network explicitly.")
			logger.LogInfo("")
			availableNetworks := utils.GetAllSupportedNetworks()

			logger.LogInfo("List of available networks: " + availableNetworks)
			logger.LogInfo("Example usage: `nodevin stop <network>`")
			logger.LogInfo("Example usage: `nodevin stop <network> --testnet`")
			logger.LogInfo("Example usage: `nodevin stop <network> --network=\"goerli\"`")
			return
		}

		network := args[0]
		if network == "all" {
			stopAllNodes()
		} else {
			stopNode(network)
		}
	},
}

func stopNode(network string) {
	logger.LogInfo("Stopping blockchain node...")

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

	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker Compose services: " + err.Error())
	}
}

func stopAllNodes() {
	logger.LogInfo("Stopping all Docker containers...")

	stopCmd := exec.Command("sh", "-c", "docker stop $(docker ps -aq)")
	stopCmd.Stdout = os.Stdout
	stopCmd.Stderr = os.Stderr

	if err := stopCmd.Run(); err != nil {
		logger.LogError("Failed to stop all Docker containers: " + err.Error())
		return
	}

	rmCmd := exec.Command("sh", "-c", "docker rm $(docker ps -aq)")
	rmCmd.Stdout = os.Stdout
	rmCmd.Stderr = os.Stderr

	if err := rmCmd.Run(); err != nil {
		logger.LogError("Failed to remove all Docker containers: " + err.Error())
		return
	}

	logger.LogInfo("All Docker containers stopped and removed successfully.")
}
