package cmd

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/docker/docker/api/types"
)

func pullImage(imageName string) error {
    ctx := context.Background()
    out, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
    if err != nil {
        return err
    }
    defer out.Close()
    fmt.Println("Pulling image:", imageName)
    return nil
}

var imageCmd = &cobra.Command{
    Use:   "image [pull]",
    Short: "Manage Docker images",
}

var pullCmd = &cobra.Command{
    Use:   "pull [image name]",
    Short: "Pull a Docker image",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        imageName := args[0]
        err := initDockerClient()
        if err != nil {
            logError("Failed to initialize Docker client: " + err.Error())
            return
        }
        err = pullImage(imageName)
        if err != nil {
            logError("Failed to pull image: " + err.Error())
            return
        }
        logInfo("Successfully pulled image: " + imageName)
    },
}

func init() {
    rootCmd.AddCommand(imageCmd)
    imageCmd.AddCommand(pullCmd)
}

