package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopNodeCmd = &cobra.Command{
	Use:   "stop-node",
	Short: "Stop a blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		stopNode()
	},
}

func stopNode() {
	logger.LogInfo("Stopping blockchain node...")

	network := viper.GetString("network")
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

func init() {
	stopNodeCmd.Flags().String("network", "bitcoin", "Blockchain network to connect to")
	viper.BindPFlag("network", stopNodeCmd.Flags().Lookup("network"))
}

/*
func stopNode() {
	logger.LogInfo("Stopping blockchain node...")

	network := viper.GetString("network")
	blockchain, exists := GetBlockchain(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	if err := blockchain.StopNode(); err != nil {
		logger.LogError("Failed to stop blockchain node: " + err.Error())
	}
}
*/
