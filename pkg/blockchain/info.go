package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/docker/docker/client"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/docker/docker/api/types/volume"
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
	args := []string{"ps", "--format", "{{json .}}"}

	// Execute the docker ps command
	cmd := exec.Command("docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.LogError("Failed to fetch Docker container information: " + err.Error())
		return
	}

	fmt.Println("Running Nodes:\n")

	// Parse the output
	containers := strings.Split(string(output), "\n")
	if len(containers) < 2 {
		fmt.Println("No running blockchain nodes found.\n")
		displayVolumeInfo()

		fmt.Println("\nHelpful Commands:\n")
		fmt.Println("nodevin delete <volume-name>\n")
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

		fmt.Fprintf(w, "| %s\t %s\t %s\t %s\t %s\n",
			container.Names,
			version,
			container.Command,
			container.Status,
			formattedPorts,
		)
	}
	w.Flush()

	fmt.Println("")

	displayVolumeInfo()

	fmt.Println("\nHelpful Commands:\n")

	fmt.Println("nodevin stop <network>")
	fmt.Println("nodevin shell <network>")
	fmt.Println("nodevin logs <network> --tail 50\n")

}

func displayVolumeInfo() {
	fmt.Println("Volume Data:\n")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.LogError("Failed to create Docker client: " + err.Error())
		return
	}

	volumeList, err := cli.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		logger.LogError("Failed to list Docker volumes: " + err.Error())
		return
	}

	if len(volumeList.Volumes) < 1 {
		fmt.Println("No blockchain node data found.")
		return
	}

	// Set up tabwriter for nicely formatted output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "| VOLUME NAME\t SIZE\t MOUNTPOINT")

	for _, vol := range volumeList.Volumes {
		size, err := getVolumeSize(vol.Mountpoint)
		sizeDescription := "unknown"
		if err == nil {
			sizeDescription = getSizeDescription(size)
		} else {
			if err.Error() == "exit status 1" {
				fmt.Println("Failed to get size for volume " + vol.Name + ", do you have proper permissions?")
			} else {
				logger.LogError("Failed to get size for volume " + vol.Name + ": " + err.Error())
			}
		}

		fmt.Fprintf(w, "| %s\t %s\t %s\n",
			vol.Name,
			sizeDescription,
			vol.Mountpoint,
		)
	}
	w.Flush()
}

func getVolumeSize(mountpoint string) (int64, error) {
	cmd := exec.Command("du", "-sb", mountpoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Parse the output to get the size in bytes
	fields := strings.Fields(string(output))
	if len(fields) == 0 {
		return 0, fmt.Errorf("unexpected du output: %s", string(output))
	}

	size, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
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

func getSizeDescription(size int64) string {
	if size <= 0 {
		return "unknown (do you have proper permissions?)"
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
