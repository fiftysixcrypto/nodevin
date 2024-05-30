package blockchain

import (
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var follow bool
var tail string

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Fetch logs from a running blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		fetchLogs()
	},
}

func fetchLogs() {
	logger.LogInfo("Fetching logs for blockchain node...")

	containerName, exists := getOutputLogsContainerName(viper.GetString("network"))
	if !exists {
		logger.LogError("Unsupported blockchain network: " + viper.GetString("network"))
		return
	}

	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	if tail != "" {
		args = append(args, "--tail", tail)
	}
	args = append(args, containerName)

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to fetch Docker logs: " + err.Error())
	}
}

func getOutputLogsContainerName(network string) (string, bool) {
	for net, container := range networkContainerMap {
		if net == network {
			return container, true
		}
	}
	return "", false
}

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVar(&tail, "tail", "all", "Number of lines to show from the end of the logs")
}
