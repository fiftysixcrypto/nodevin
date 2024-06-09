package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/internal/utils"
	"github.com/curveballdaniel/nodevin/pkg/blockchain/bitcoin"
	"github.com/curveballdaniel/nodevin/pkg/docker"

	"github.com/spf13/cobra"
)

var ord bool

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
		logger.LogError("No network provided. Nodevin supports any of the following: " + utils.GetAllSupportedNetworks())
		logger.LogInfo("Example usage: `nodevin start <network>`")
		return
	}

	network := args[0]

	containerName, exists := utils.GetFiftysixDockerhubContainerName(network)
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

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logger.LogError("failed to get current working directory: " + err.Error())
		return
	}

	// Create env file for chain compose
	composeFilePath, err := createComposeFileForNetwork(network, cwd)
	if err != nil {
		logger.LogError("Failed to create node docker compose file: " + err.Error())
	}

	// Start the node
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to start Docker Compose services: " + err.Error())
		return
	}

	logger.LogInfo("Successfully started blockchain node for network: " + network)

	startMessage, _ := utils.GetStartMessage(network)

	fmt.Printf("\n%s\n", startMessage)
}

func createComposeFileForNetwork(network string, cwd string) (string, error) {
	switch network {
	case "bitcoin":
		return bitcoin.CreateBitcoinComposeFile(cwd)
	case "ethereum":
		return "", nil // createEthereumComposeFile(cwd)
	case "dogecoin":
		return "", nil // createDogecoinComposeFile(cwd)
	case "ethereumclassic":
		return "", nil // createEthereumClassicComposeFile(cwd)
	case "litecoin":
		return "", nil // createLitecoinComposeFile(cwd)
	default:
		return "", fmt.Errorf("unsupported network: %s", network)
	}
}
