package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/update"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the NodeVin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		runDaemon()
	},
}

func runDaemon() {
	logger.LogInfo("Starting NodeVin daemon...")

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Set up a ticker to check for updates periodically
	updateTicker := time.NewTicker(24 * time.Hour)
	imageCheckTicker := time.NewTicker(1 * time.Hour)

	go func() {
		for {
			select {
			case <-updateTicker.C:
				logger.LogInfo("Checking for updates...")
				updateNeeded, err := update.CheckForUpdates()
				if err != nil {
					logger.LogError("Failed to check for updates: " + err.Error())
					continue
				}
				if updateNeeded {
					logger.LogInfo("Update downloaded. Applying update...")
					if err := update.ApplyUpdate(); err != nil {
						logger.LogError("Failed to apply update: " + err.Error())
						continue
					}
					logger.LogInfo("Update applied successfully. Please restart the application.")
				} else {
					logger.LogInfo("Nodevin is up to date.")
				}
			case <-imageCheckTicker.C:
				logger.LogInfo("Checking for Docker image updates...")
				if err := checkAndUpdateDockerImages(); err != nil {
					logger.LogError("Failed to check/update Docker images: " + err.Error())
				}
			}
		}
	}()

	// Block until a signal is received
	sig := <-sigs
	logger.LogInfo("Received signal: " + sig.String())

	// Perform cleanup and shutdown
	shutdownDaemon()
	logger.LogInfo("Shutting down NodeVin daemon...")
}

func shutdownDaemon() {
	// Implement the cleanup and shutdown logic here
}

func Execute() error {
	return Cmd.Execute()
}

func checkAndUpdateDockerImages() error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})

	if err != nil {
		return fmt.Errorf("failed to list Docker containers: %w", err)
	}

	for _, container := range containers {
		if strings.HasPrefix(container.Image, "fiftysix/") {
			image := strings.Split(container.Image, ":")[0]
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
				logger.LogInfo(fmt.Sprintf("Updating image %s to latest version", container.Image))
				if err := updateDockerImage(cli, container, image); err != nil {
					logger.LogError(fmt.Sprintf("Failed to update image %s: %v", container.Image, err))
				}
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
	return inspect.RepoDigests[0], nil
}

func updateDockerImage(cli *client.Client, container types.Container, image string) error {
	composeFilePath := fmt.Sprintf("docker-compose_%s.yml", container.Names[0])

	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Docker Compose services: %w", err)
	}

	pullCmd := exec.Command("docker", "pull", image+":latest")
	if err := pullCmd.Run(); err != nil {
		return fmt.Errorf("failed to pull latest Docker image: %w", err)
	}

	cmd = exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker Compose services: %w", err)
	}

	return nil
}
