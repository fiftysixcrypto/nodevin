package nodes

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

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

	// Get Docker images
	listCmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}")
	var out bytes.Buffer
	listCmd.Stdout = &out
	listCmd.Stderr = os.Stderr

	if err := listCmd.Run(); err != nil {
		logger.LogError("Failed to list Docker images: " + err.Error())
		return
	}

	// Filter images that start with 'fiftysix/'
	images := strings.Split(out.String(), "\n")
	var filteredImages []string
	for _, image := range images {
		if strings.HasPrefix(image, "fiftysix/") {
			filteredImages = append(filteredImages, image)
		}
	}

	// Remove filtered images
	if len(filteredImages) > 0 {
		logger.LogInfo("Removing images: " + strings.Join(filteredImages, ", "))
		removeCmd := exec.Command("docker", append([]string{"rmi"}, filteredImages...)...)
		removeCmd.Stdout = os.Stdout
		removeCmd.Stderr = os.Stderr
		if err := removeCmd.Run(); err != nil {
			logger.LogError("Failed to remove Docker images: " + err.Error())
			return
		}
	} else {
		logger.LogInfo("No Docker images starting with 'fiftysix/' found.")
	}

	logger.LogInfo("Successfully removed all Docker images starting with 'fiftysix/'.")
}
