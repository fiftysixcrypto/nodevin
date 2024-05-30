package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/docker"
	"github.com/spf13/cobra"
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
		displayInfo()
	},
}

func displayInfo() {
	// Prepare the docker ps command
	args := []string{"ps", "-a", "--format", "{{json .}}"}

	// Execute the docker ps command
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogError("Failed to fetch Docker container information: " + err.Error())
		return
	}

	// Parse the output
	containers := strings.Split(string(output), "\n")
	if len(containers) < 2 {
		logger.LogInfo("No running blockchain nodes found.")
		return
	}

	// Set up tabwriter for nicely formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "| BLOCKCHAIN\t VERSION\t COMMAND\t STATUS\t PORTS")

	for _, containerJSON := range containers {
		if strings.TrimSpace(containerJSON) == "" {
			continue
		}
		var container ContainerInfo
		if err := json.Unmarshal([]byte(containerJSON), &container); err != nil {
			logger.LogError("Failed to parse container JSON: " + err.Error())
			continue
		}

		imageName := container.Image
		if strings.HasPrefix(container.Image, "fiftysix/") {
			imageName = strings.TrimPrefix(container.Image, "fiftysix/")
		} else {
			// Skip if the image does not begin with fiftysix
			continue
		}

		version := "unknown"
		if parts := strings.Split(imageName, ":"); len(parts) > 1 {
			imageName = parts[0]
			version = parts[1]
		}

		formattedPorts := formatPorts(container.Ports)

		volumeInfo, err := docker.ListVolumeDetails(imageName)
		if err != nil {
			logger.LogError("Failed to fetch Docker volume information: " + err.Error())
		}

		stopCmd := fmt.Sprintf("nodevin stop-node --network %s", getNetworkFromImage(container.Image))
		logsCmd := fmt.Sprintf("nodevin logs --network %s", getNetworkFromImage(container.Image))

		fmt.Fprintf(w, "| %s\t %s\t %s\t %s\t %s\n\n%s\n%s\n\n%s\n%s\n",
			imageName,
			version,
			container.Command,
			container.Status,
			formattedPorts,
			"Data Location: "+volumeInfo.Mountpoint,
			"Data Size: "+getSizeDescription(int64(volumeInfo.Size)),
			"Node Logs: "+logsCmd,
			"Stop Node: "+stopCmd,
		)
	}
	w.Flush()
}

func formatPorts(ports string) string {
	// Split the ports field by commas
	portSegments := strings.Split(ports, ",")
	formattedPorts := []string{}
	uniquePorts := make(map[string]bool)

	for _, segment := range portSegments {
		// Use a regex to extract the port numbers and ranges
		re := regexp.MustCompile(`(\d+(-\d+)?)`)
		matches := re.FindAllString(segment, -1)
		for _, match := range matches {
			if match != "0" && !uniquePorts[match] {
				uniquePorts[match] = true
				formattedPorts = append(formattedPorts, match)
			}
		}
	}
	return strings.Join(formattedPorts, ", ")
}

func getNetworkFromImage(image string) string {
	for network, container := range networkContainerMap {
		if strings.Contains(image, container) {
			return network
		}
	}
	return "unknown"
}

func getSizeDescription(size int64) string {
	if size <= 0 {
		return "unknown (do you have permissions?)"
	}

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
		PB = TB * 1024
	)

	switch {
	case size >= PB:
		return fmt.Sprintf("%.2f PB", float64(size)/PB)
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}
