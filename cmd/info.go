package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ContainerInfo struct {
	ID         string `json:"ID"`
	Image      string `json:"Image"`
	Command    string `json:"Command"`
	CreatedAt  string `json:"CreatedAt"`
	RunningFor string `json:"RunningFor"`
	Status     string `json:"Status"`
	Ports      string `json:"Ports"`
	Names      string `json:"Names"`
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get information about all running blockchain nodes",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		displayInfo()
	},
}

func displayInfo() {
	logInfo("Fetching information about running blockchain nodes...")

	// Prepare the docker ps command
	args := []string{"ps", "-a", "--format", "{{json .}}"}

	// Execute the docker ps command
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logError("Failed to fetch Docker container information: " + err.Error())
		return
	}

	// Parse the output
	containers := strings.Split(string(output), "\n")
	if len(containers) < 2 {
		logInfo("No running blockchain nodes found.")
		return
	}

	// Set up tabwriter for nicely formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\nBLOCKCHAIN\tVERSION\tCOMMAND\tSTATUS\tPORTS")

	for _, containerJSON := range containers {
		if strings.TrimSpace(containerJSON) == "" {
			continue
		}
		var container ContainerInfo
		if err := json.Unmarshal([]byte(containerJSON), &container); err != nil {
			logError("Failed to parse container JSON: " + err.Error())
			continue
		}

		_, exists := getDockerContainerName(container.Image)
		if !exists {
			logError("Unsupported blockchain network for image: " + container.Image)
			continue
		}

		version := "unknown"
		if parts := strings.Split(container.Image, ":"); len(parts) > 1 {
			version = parts[1]
		}

		stopCmd := fmt.Sprintf("go run main.go stop-node --network %s", getNetworkFromImage(container.Image))
		logsCmd := fmt.Sprintf("go run main.go logs --network %s", getNetworkFromImage(container.Image))

		// TODO: container.Image and version should be smooshed together to create bitcoin-core:21.0

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n\n%s\n%s\n",
			container.Image,
			version,
			container.Command,
			container.Status,
			container.Ports,
			"Logs Command: "+logsCmd,
			"Stop Command: "+stopCmd,
		)
	}
	w.Flush()
}

func getDockerContainerName(image string) (string, bool) {
	for _, container := range networkContainerMap {
		if strings.Contains(image, container) {
			return container, true
		}
	}
	return "", false
}

func getNetworkFromImage(image string) string {
	for network, container := range networkContainerMap {
		if strings.Contains(image, container) {
			return network
		}
	}
	return "unknown"
}

func init() {
	infoCmd.Flags().StringVar(&config.Network, "network", "mainnet", "Blockchain network to connect to")

	viper.BindPFlag("network", infoCmd.Flags().Lookup("network"))

	rootCmd.AddCommand(infoCmd)
}
