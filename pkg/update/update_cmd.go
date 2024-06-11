package update

import (
	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var (
	UpdateCmd = updateCmd
)

var updateCmd = &cobra.Command{
	Use:   "update [target]",
	Short: "Update NodeVin software or Docker images",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Update NodeVin software
			logger.LogInfo("Checking for Nodevin updates...")
			updateNeeded, err := CheckForUpdates()
			if err != nil {
				logger.LogError("Failed to check for updates: " + err.Error())
				return
			}
			if updateNeeded {
				logger.LogInfo("Update downloaded. Applying update...")
				if err := ApplyUpdate(); err != nil {
					logger.LogError("Failed to apply update: " + err.Error())
					return
				}
				logger.LogInfo("Update applied successfully. Please restart the application.")
			} else {
				logger.LogInfo("Nodevin is up to date.")
			}
		} else if args[0] == "docker" {
			// Update Docker images
			logger.LogInfo("Checking for Docker image updates...")
			if err := CheckAndUpdateDockerImages(); err != nil {
				logger.LogError("Failed to check/update Docker images: " + err.Error())
				return
			}
			logger.LogInfo("Docker images updated successfully.")
		} else {
			logger.LogError("Invalid target specified. Use 'nodevin update' or 'nodevin update docker'.")
		}
	},
}
