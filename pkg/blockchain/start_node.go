package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/blockchain/bitcoin"
	"github.com/curveballdaniel/nodevin/pkg/docker"
	"github.com/spf13/cobra"
)

var startNodeCmd = &cobra.Command{
	Use:   "start [network]",
	Short: "Start a blockchain node",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		startNode(args)
	},
}

func startNode(args []string) {
	if len(args) == 0 {
		logger.LogError("No network provided. Nodevin supports any of the following: " + getAllSupportedNetworks())
		logger.LogInfo("Example usage: `nodevin start <network>`")
		return
	}

	network := args[0]

	containerName, exists := getFiftysixDockerhubContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	logger.LogInfo("Starting blockchain node for network: " + network)

	// Initialize Docker client
	if err := docker.InitDockerClient(); err != nil {
		logger.LogError("Failed to initialize Docker client: " + err.Error())
		return
	}

	// Pull the Docker image
	image := containerName + ":latest"
	if err := docker.PullImage(image); err != nil {
		logger.LogError("Failed to pull Docker image: " + err.Error())
		return
	}

	// Obtain compose name from mappings
	dockerContainerName, exists := getFiftysixLocalMappedContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain compose file: " + network)
		return
	}

	fmt.Println("here's the chain name", dockerContainerName, containerName, network)

	// Set env variables for chain compose
	envFilePath, err := bitcoin.CreateBitcoinEnv()
	if err != nil {
		logger.LogError("Failed to create node settings .env file: " + err.Error())
	}

	// Start the node
	composeFilePath := fmt.Sprintf("docker/%s/docker-compose_%s.yml", network, dockerContainerName)
	cmd := exec.Command("docker-compose", "--env-file", envFilePath, "-f", composeFilePath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to start Docker Compose services: " + err.Error())
		return
	}

	logger.LogInfo("Successfully started blockchain node for network: " + network)
}
