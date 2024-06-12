package update

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CheckAndUpdateDockerImages() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list Docker containers: %w", err)
	}

	if len(containers) == 0 {
		logger.LogInfo("No Docker containers running!")
		return nil
	}

	for _, container := range containers {
		if strings.HasPrefix(container.Image, "fiftysix/") {
			image := strings.TrimPrefix(strings.Split(container.Image, ":")[0], "fiftysix/")
			latestDigest, err := getDockerHubImageDigest("fiftysix", image, "latest")
			if err != nil {
				logger.LogError(fmt.Sprintf("Failed to get latest digest for image %s: %v", image, err))
				continue
			}

			localDigest, err := getLocalImageDigest(cli, container.Image)
			if err != nil {
				logger.LogError(fmt.Sprintf("Failed to get local digest for image %s: %v", container.Image, err))
				continue
			}

			if latestDigest != localDigest {
				logger.LogInfo(fmt.Sprintf("New version of %s found!", image))
				logger.LogInfo(fmt.Sprintf("Attempting to update %s to latest version", image))
				if err := updateDockerImage(container, container.Image); err != nil {
					logger.LogError(fmt.Sprintf("Failed to update image %s: %v", image, err))
				}
			} else {
				logger.LogInfo(fmt.Sprintf("Local image %s is on latest version", image))
			}
		}
	}

	return nil
}

func getDockerHubImageDigest(namespace, repository, tag string) (string, error) {
	url := fmt.Sprintf("https://hub.docker.com/v2/namespaces/%s/repositories/%s/tags/%s", namespace, repository, tag)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tags: %s", resp.Status)
	}

	var result struct {
		Images []struct {
			Digest string `json:"digest"`
		} `json:"images"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Images) == 0 {
		return "", fmt.Errorf("no images found for tag %s", tag)
	}

	return result.Images[0].Digest, nil
}

func getLocalImageDigest(cli *client.Client, image string) (string, error) {
	inspect, _, err := cli.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		return "", err
	}
	if len(inspect.RepoDigests) == 0 {
		return "", fmt.Errorf("no repo digests found for image %s", image)
	}
	localDigest := inspect.RepoDigests[0]
	if parts := strings.Split(localDigest, "@"); len(parts) > 1 {
		return parts[1], nil
	}
	return localDigest, nil
}

func updateDockerImage(container types.Container, image string) error {
	imageShorthandName := strings.TrimPrefix(container.Names[0], "/")
	composeFilePath := fmt.Sprintf("docker-compose_%s.yml", imageShorthandName)

	logger.LogInfo(fmt.Sprintf("Shutting down %s...", imageShorthandName))
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Docker Compose services: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("Pulling latest image for %s...", imageShorthandName))
	pullCmd := exec.Command("docker", "pull", image)
	if err := pullCmd.Run(); err != nil {
		return fmt.Errorf("failed to pull latest Docker image: %w", err)
	}
	logger.LogInfo(fmt.Sprintf("Successfully pulled latest image for %s", imageShorthandName))

	logger.LogInfo(fmt.Sprintf("Starting %s back up on latest version...", imageShorthandName))
	cmd = exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker Compose services: %w", err)
	}

	return nil
}