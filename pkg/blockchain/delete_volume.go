package blockchain

import (
	"os"
	"os/exec"
	"strings"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/docker"
	"github.com/spf13/cobra"
)

var deleteVolumeCmd = &cobra.Command{
	Use:   "delete [volume-name]",
	Short: "Delete a Docker volume",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.LogError("No volume name provided. To delete a volume, specify the volume name explicitly.")
			dockerVolumes, err := docker.ListVolumes()
			if err != nil {
				logger.LogError("Unable to list local Docker volumes: " + err.Error())
				return
			}
			logger.LogInfo("List of current volumes: " + dockerVolumes)
			logger.LogInfo("Example usage: `nodevin delete <volume-name>`")
			logger.LogInfo("Example usage: `nodevin delete all`")
			return
		}

		volumeName := args[0]

		err := docker.InitDockerClient()
		if err != nil {
			logger.LogError("Failed to initialize Docker client: " + err.Error())
			return
		}

		if volumeName == "all" {
			deleteAllVolumes()
		} else {
			deleteVolume(volumeName)
		}
	},
}

func deleteVolume(volumeName string) {
	volumes, err := docker.ListVolumes()
	if err != nil {
		logger.LogError("Unable to list local Docker volumes: " + err.Error())
		return
	}

	if !strings.Contains(volumes, volumeName) {
		logger.LogError("Volume name not found: " + volumeName)
		logger.LogInfo("List of current volumes: " + volumes)
		logger.LogInfo("Example usage: `nodevin delete <volume-name>`")
		logger.LogInfo("Example usage: `nodevin delete all`")
		return
	}

	completed, err := docker.RemoveVolume(volumeName)
	if err != nil {
		logger.LogError("Failed to remove volume: " + err.Error())
		return
	}

	if completed {
		logger.LogInfo("Successfully removed volume: " + volumeName)
	} else {
		logger.LogError("Failed to remove volume: " + volumeName)
	}
}

func deleteAllVolumes() {
	logger.LogInfo("Deleting all Docker volumes with label 'nodevin.blockchain.software'...")

	// Run the Docker command to remove all volumes with the nodevin volume label
	removeCmd := exec.Command("sh", "-c", "docker volume ls -q -f label=nodevin.blockchain.software | xargs -r docker volume rm")
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr

	if err := removeCmd.Run(); err != nil {
		logger.LogError("Failed to remove Docker volumes with label 'nodevin.blockchain.software': " + err.Error())
		return
	}

	logger.LogInfo("All Docker volumes with label 'nodevin.blockchain.software' removed successfully.")
}
