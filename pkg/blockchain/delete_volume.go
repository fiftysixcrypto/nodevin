package blockchain

import (
	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/internal/utils"
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
				logger.LogError("Unable to list local Docker volumes.")
				return
			}
			logger.LogInfo("List of current volumes: " + dockerVolumes)
			logger.LogInfo("Example usage: `nodevin delete <volume-name>`")
			return
		}

		passedVolumeArg := args[0]

		err := docker.InitDockerClient()
		if err != nil {
			logger.LogError("Failed to initialize Docker client: " + err.Error())
			return
		}
		completed, err := docker.RemoveVolume(passedVolumeArg)

		if err != nil {
			logger.LogError("Failed to remove volume: " + err.Error())
			return
		}

		if !completed {
			_, exists := utils.GetFiftysixDockerhubContainerName(passedVolumeArg)

			if !exists {
				logger.LogError("To delete a volume list the volume name explicitly.")
				dockerVolumes, err := docker.ListVolumes()
				if err != nil {
					logger.LogError("Unable to list local Docker volumes.")
					return
				}
				logger.LogInfo("List of current volumes: " + dockerVolumes)
				logger.LogInfo("Example usage: `nodevin delete <volume-name>`")
				return
			}

			return
		}

		logger.LogInfo("Successfully removed volume: " + passedVolumeArg)
	},
}
