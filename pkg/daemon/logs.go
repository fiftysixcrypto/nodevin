package daemon

import (
	"os"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var daemonLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show logs of the nodevin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		showDaemonLogs()
	},
}

func showDaemonLogs() {
	logData, err := os.ReadFile(logFilePath)
	if err != nil {
		logger.LogError("Failed to read log file: " + err.Error())
		return
	}

	logger.LogInfo(string(logData))
}
