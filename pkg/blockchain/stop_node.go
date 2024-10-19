package blockchain

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/utils"
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
			availableNetworks := utils.GetCommandSupportedNetworks()

			logger.LogInfo("List of available networks: " + availableNetworks)
			logger.LogInfo(fmt.Sprintf("Example usage: `%s stop <network>`", utils.GetNodevinExecutable()))
			logger.LogInfo(fmt.Sprintf("Example usage: `%s stop <network> --testnet`", utils.GetNodevinExecutable()))
			logger.LogInfo(fmt.Sprintf("Example usage: `%s stop all`", utils.GetNodevinExecutable()))
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

	containerName, exists := utils.GetDefaultLocalMappedContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	// Get the user's home directory in a cross-platform manner
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to determine home directory: %s", err))
		return
	}

	var composeCreateDir string

	// Check if ~/.nodevin exists
	nodevinDir := filepath.Join(homeDir, ".nodevin", "data")
	if _, err := os.Stat(nodevinDir); err == nil {
		composeCreateDir = nodevinDir
	} else if !os.IsNotExist(err) {
		logger.LogError(fmt.Sprintf("Error accessing ~/.nodevin: %s", err))
	}

	// Fallback to nodevin executable directory if ~/.nodevin does not exist
	if composeCreateDir == "" {
		composeCreatePath, err := os.Executable()
		if err != nil {
			cwd, wdErr := os.Getwd()
			if wdErr != nil {
				logger.LogError("Unable to determine executable or working directory")
				return
			}
			composeCreateDir = cwd
		} else {
			// Use the directory where the executable is located
			composeCreateDir = filepath.Dir(composeCreatePath)
		}
	}

	composeFileName := fmt.Sprintf("docker-compose_%s.yml", containerName)
	if utils.CheckIfTestnetOrTestnetNetworkFlag() {
		composeFileName = fmt.Sprintf("docker-compose_%s.yml", containerName+"-testnet")
	}

	composeFilePath := filepath.Join(composeCreateDir, composeFileName)

	// Check if there are any running containers for this compose file
	psCmd := exec.Command("docker-compose", "-f", composeFilePath, "ps", "-q")
	psOut, err := psCmd.Output()
	if err != nil {
		logger.LogError("Failed to find Docker Compose services: " + err.Error())
		return
	}

	if len(psOut) == 0 {
		logger.LogInfo("No running containers found for the specified network (did you mean to add --testnet?)")
		return
	}

	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker Compose services: " + err.Error())
		return
	}

	logger.LogInfo("Blockchain node stopped successfully.")
}

func stopAllNodes() {
	logger.LogInfo("Stopping all Docker containers...")

	// Get a list of running Docker container IDs
	psCmd := exec.Command("docker", "ps", "-q")
	var psOut bytes.Buffer
	psCmd.Stdout = &psOut
	psCmd.Stderr = os.Stderr

	if err := psCmd.Run(); err != nil {
		logger.LogError("Failed to list running Docker containers: " + err.Error())
		return
	}

	containerIDs := strings.Fields(psOut.String())
	if len(containerIDs) == 0 {
		logger.LogInfo("No running Docker containers found.")
		return
	}

	// Stop the containers
	logger.LogInfo("Stopping containers: " + strings.Join(containerIDs, ", "))
	stopCmd := exec.Command("docker", append([]string{"stop"}, containerIDs...)...)
	stopCmd.Stdout = os.Stdout
	stopCmd.Stderr = os.Stderr

	if err := stopCmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker containers: " + err.Error())
		return
	}

	// Remove the containers
	logger.LogInfo("Removing containers: " + strings.Join(containerIDs, ", "))
	rmCmd := exec.Command("docker", append([]string{"rm"}, containerIDs...)...)
	rmCmd.Stdout = os.Stdout
	rmCmd.Stderr = os.Stderr

	if err := rmCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker containers: " + err.Error())
		return
	}

	logger.LogInfo("All Docker containers stopped and removed successfully.")
}
