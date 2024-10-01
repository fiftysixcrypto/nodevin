package blockchain

import (
	"os"
	"os/exec"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup all Docker images starting with 'fiftysix/'",
	Run: func(cmd *cobra.Command, args []string) {
		cleanupAllImages()
	},
}

func cleanupAllImages() {
	logger.LogInfo("Removing all Docker images starting with 'fiftysix/'...")

	removeCmd := exec.Command("sh", "-c", "docker images --format '{{.Repository}}:{{.Tag}}' | grep '^fiftysix/' | xargs -r docker rmi")
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker images: " + err.Error())
		return
	}

	logger.LogInfo("Successfully removed all Docker images starting with 'fiftysix/'.")
}
