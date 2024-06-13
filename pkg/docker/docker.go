package docker

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/fiftysixcrypto/nodevin/internal/logger"
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

type VolumeDetails struct {
	Name        string
	Size        int
	Mountpoint  string
	DateCreated string
}

func ListVolumeDetails(network string) (*VolumeDetails, error) {
	if dockerClient == nil {
		if err := InitDockerClient(); err != nil {
			return nil, fmt.Errorf("failed to initialize Docker client: %w", err)
		}
	}

	ctx := context.Background()
	volumes, err := dockerClient.VolumeList(ctx, volume.ListOptions{
		Filters: filters.NewArgs(filters.Arg("label", fmt.Sprintf("nodevin.blockchain.software=%s", network))),
	})

	if err != nil {
		return nil, err
	}

	if len(volumes.Volumes) == 0 {
		return nil, fmt.Errorf("no volumes found for network: %s", network)
	}

	vol := volumes.Volumes[0]

	volInspect, err := dockerClient.VolumeInspect(ctx, vol.Name)
	if err != nil {
		return nil, err
	}

	size := -1

	if volInspect.Mountpoint != "" {
		obtainedSize, _ := CalculateDirSize(volInspect.Mountpoint)
		if int(obtainedSize) > 0 {
			size = int(obtainedSize)
		}
	}

	return &VolumeDetails{
		Name:        vol.Name,
		Mountpoint:  vol.Mountpoint,
		DateCreated: volInspect.CreatedAt,
		Size:        size,
	}, nil
}

func CalculateDirSize(path string) (int64, error) {
	cmd := exec.Command("du", "-sb", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return 0, err
	}
	var size int64
	_, err := fmt.Sscanf(out.String(), "%d", &size)
	return size, err
}

func ListVolumes() (string, error) {
	if dockerClient == nil {
		if err := InitDockerClient(); err != nil {
			return "", fmt.Errorf("failed to initialize Docker client: %w", err)
		}
	}

	ctx := context.Background()

	volumes, err := dockerClient.VolumeList(ctx, volume.ListOptions{
		Filters: filters.Args{},
	})
	if err != nil {
		return "", err
	}

	var volumeNames []string
	for _, v := range volumes.Volumes {
		volumeNames = append(volumeNames, v.Name)
	}

	finalResponse := strings.Join(volumeNames, ", ")
	if len(finalResponse) < 1 {
		finalResponse = "(none)"
	}

	return finalResponse, nil
}

func RemoveVolume(volumeName string) (bool, error) {
	ctx := context.Background()

	// Check if the volume exists
	_, err := dockerClient.VolumeInspect(ctx, volumeName)
	if err != nil {
		if client.IsErrNotFound(err) {
			// Volume does not exist
			return false, nil
		}
		// An error occurred while inspecting the volume
		return false, err
	}

	// Volume exists, attempt to remove it
	err = dockerClient.VolumeRemove(ctx, volumeName, true)
	if err != nil {
		return true, err // true because it found the volume, but there was an error
	}

	fmt.Println("Removed volume:", volumeName)
	return true, nil
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
