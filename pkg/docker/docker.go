package docker

import (
	"context"
	"fmt"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

var dockerClient *client.Client

func InitDockerClient() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.LogError("Failed to initialize Docker client: " + err.Error())
		return err
	}
	dockerClient = cli
	return nil
}

func CreateVolume(volumeName string) error {
	ctx := context.Background()
	_, err := dockerClient.VolumeCreate(ctx, volume.CreateOptions{
		Name: volumeName,
	})
	if err != nil {
		return err
	}
	fmt.Println("Created volume:", volumeName)
	return nil
}

func ListVolumes() error {
	ctx := context.Background()
	volumes, err := dockerClient.VolumeList(ctx, volume.ListOptions{
		Filters: filters.Args{},
	})
	if err != nil {
		return err
	}
	for _, v := range volumes.Volumes {
		fmt.Println("Volume:", v.Name)
	}
	return nil
}

func RemoveVolume(volumeName string) error {
	ctx := context.Background()
	err := dockerClient.VolumeRemove(ctx, volumeName, true)
	if err != nil {
		return err
	}
	fmt.Println("Removed volume:", volumeName)
	return nil
}

func PullImage(image string) error {
	logger.LogInfo("Pulling Docker image: " + image)
	out, err := dockerClient.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		logger.LogError("Failed to pull Docker image: " + err.Error())
		return err
	}
	defer out.Close()
	fmt.Println(out)
	return nil
}
