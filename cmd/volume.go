package cmd

import (
    "github.com/spf13/cobra"
)

var volumeCmd = &cobra.Command{
    Use:   "volume [create|list|remove]",
    Short: "Manage Docker volumes",
}

var createVolumeCmd = &cobra.Command{
    Use:   "create [volume name]",
    Short: "Create a Docker volume",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        volumeName := args[0]
        err := initDockerClient()
        if err != nil {
            logError("Failed to initialize Docker client: " + err.Error())
            return
        }
        err = createVolume(volumeName)
        if err != nil {
            logError("Failed to create volume: " + err.Error())
            return
        }
        logInfo("Successfully created volume: " + volumeName)
    },
}

var listVolumeCmd = &cobra.Command{
    Use:   "list",
    Short: "List all Docker volumes",
    Run: func(cmd *cobra.Command, args []string) {
        err := initDockerClient()
        if err != nil {
            logError("Failed to initialize Docker client: " + err.Error())
            return
        }
        err = listVolumes()
        if err != nil {
            logError("Failed to list volumes: " + err.Error())
            return
        }
    },
}

var removeVolumeCmd = &cobra.Command{
    Use:   "remove [volume name]",
    Short: "Remove a Docker volume",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        volumeName := args[0]
        err := initDockerClient()
        if err != nil {
            logError("Failed to initialize Docker client: " + err.Error())
            return
        }
        err = removeVolume(volumeName)
        if err != nil {
            logError("Failed to remove volume: " + err.Error())
            return
        }
        logInfo("Successfully removed volume: " + volumeName)
    },
}

func init() {
    rootCmd.AddCommand(volumeCmd)
    volumeCmd.AddCommand(createVolumeCmd)
    volumeCmd.AddCommand(listVolumeCmd)
    volumeCmd.AddCommand(removeVolumeCmd)
}

