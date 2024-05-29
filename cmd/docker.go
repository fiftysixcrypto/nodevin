package cmd

import (
    "context"
    "fmt"
    "github.com/docker/docker/api/types/filters"
    "github.com/docker/docker/api/types/volume"
    "github.com/docker/docker/client"
)

var dockerClient *client.Client

func initDockerClient() error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if (err != nil) {
        return err
    }
    dockerClient = cli
    return nil
}

func createVolume(volumeName string) error {
    ctx := context.Background()
    _, err := dockerClient.VolumeCreate(ctx, volume.CreateOptions{
        Name: volumeName,
    })
    if (err != nil) {
        return err
    }
    fmt.Println("Created volume:", volumeName)
    return nil
}

func listVolumes() error {
    ctx := context.Background()
    volumes, err := dockerClient.VolumeList(ctx, volume.ListOptions{
        Filters: filters.Args{},
    })
    if (err != nil) {
        return err
    }
    for _, v := range volumes.Volumes {
        fmt.Println("Volume:", v.Name)
    }
    return nil
}

func removeVolume(volumeName string) error {
    ctx := context.Background()
    err := dockerClient.VolumeRemove(ctx, volumeName, true)
    if (err != nil) {
        return err
    }
    fmt.Println("Removed volume:", volumeName)
    return nil
}

