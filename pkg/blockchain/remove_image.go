package blockchain

import (
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/internal/utils"
	"github.com/curveballdaniel/nodevin/pkg/docker"

	"github.com/spf13/cobra"
)

var removeImageCmd = &cobra.Command{
	Use:   "remove-image [network]",
	Short: "Clear Docker images for a blockchain network",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeImageImages(args)
	},
}

func removeImageImages(args []string) {
	if len(args) == 0 {
		logger.LogError("No network provided. Specify a network to remove images for, or use 'all' to remove all images starting with 'fiftysix/'.")
		logger.LogInfo("Example usage: `nodevin remove-image <network>`")
		logger.LogInfo("Example usage: `nodevin remove-image all`")
		return
	}

	network := args[0]

	err := docker.InitDockerClient()
	if err != nil {
		logger.LogError("Failed to initialize Docker client: " + err.Error())
		return
	}

	if network == "all" {
		removeImageAllImages()
	} else {
		removeImageNetworkImage(network)
	}
}

func removeImageNetworkImage(network string) {
	imageName, exists := utils.GetFiftysixDockerhubContainerName(network)
	if !exists {
		logger.LogError("Unsupported blockchain network: " + network)
		return
	}

	image := imageName + ":latest"
	logger.LogInfo("Clearing Docker image for network: " + network)

	removeCmd := exec.Command("docker", "rmi", image)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker image: " + err.Error())
		return
	}

	logger.LogInfo("Successfully removed Docker image: " + image)
}

func removeImageAllImages() {
	logger.LogInfo("Clearing all Docker images starting with 'fiftysix/'...")

	removeCmd := exec.Command("sh", "-c", "docker images --format '{{.Repository}}:{{.Tag}}' | grep '^fiftysix/' | xargs -r docker rmi")
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker images: " + err.Error())
		return
	}

	logger.LogInfo("Successfully removed all Docker images starting with 'fiftysix/'.")
}
