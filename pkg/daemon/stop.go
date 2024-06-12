package daemon

import (
	"os"
	"strconv"
	"syscall"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the nodevin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		stopDaemon()
	},
}

func stopDaemon() {
	pidData, err := os.ReadFile(pidFilePath)
	if err != nil {
		logger.LogError("Failed to read PID file: " + err.Error())
		logger.LogInfo("Are you sure the daemon is running?")
		return
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		logger.LogError("Invalid PID: " + err.Error())
		return
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		logger.LogError("Failed to find process: " + err.Error())
		return
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		logger.LogError("Failed to stop daemon: " + err.Error())
		return
	}

	err = os.Remove(pidFilePath)
	if err != nil {
		logger.LogError("Failed to remove PID file: " + err.Error())
		return
	}

	logger.LogInfo("Nodevin daemon successfully stopped.")
}
