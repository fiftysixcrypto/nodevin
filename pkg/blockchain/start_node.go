package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startNodeCmd = &cobra.Command{
	Use:   "start-node",
	Short: "Start a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		startNode()
	},
}

func startNode() {
	logger.LogInfo("Starting blockchain node...")

	network := viper.GetString("network")
	containerName, exists := getFiftysixDockerhubContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

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

	// Start the node
	dockerContainerName, exists := getFiftysixLocalMappedContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain compose file: " + network)
		return
	}

	composeFilePath := fmt.Sprintf("docker/%s/docker-compose_%s.yml", network, dockerContainerName)
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to start Docker Compose services: " + err.Error())
	}
}

func init() {
	startNodeCmd.Flags().String("network", "bitcoin", "Blockchain network to connect to")
	viper.BindPFlag("network", startNodeCmd.Flags().Lookup("network"))
}

/*
func startNode() {
	logger.LogInfo("Starting blockchain node...")

	network := viper.GetString("network")
	blockchain, exists := GetBlockchain(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	if err := blockchain.StartNode(config.AppConfig); err != nil {
		logger.LogError("Failed to start blockchain node: " + err.Error())
	}
}
*/
