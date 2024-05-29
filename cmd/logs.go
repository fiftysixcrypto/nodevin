package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	follow bool
	tail   string
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Fetch logs from a running blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		fetchLogs()
	},
}

func fetchLogs() {
	logInfo("Fetching logs for blockchain node...")

	// Log configuration for debugging
	logInfo(fmt.Sprintf("Network: %s", config.Network))

	// Get the container name from the network mapping
	containerName, exists := getContainerName(config.Network)
	if !exists {
		logError("Unsupported blockchain network: " + config.Network)
		return
	}

	// Prepare the docker logs command
	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	if tail != "" {
		args = append(args, "--tail", tail)
	}
	args = append(args, containerName)

	// Execute the docker logs command
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logError("Failed to fetch Docker logs: " + err.Error())
	}
}

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVar(&tail, "tail", "all", "Number of lines to show from the end of the logs")

	logsCmd.Flags().StringVar(&config.Network, "network", "mainnet", "Blockchain network to connect to")

	viper.BindPFlag("network", logsCmd.Flags().Lookup("network"))

	rootCmd.AddCommand(logsCmd)
}
