package daemon

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var (
	lines  int
	clear  bool
	follow bool
)

var daemonLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Manage logs of the nodevin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if clear {
			clearDaemonLogs()
		} else {
			showDaemonLogs()
		}
	},
}

func init() {
	daemonLogsCmd.Flags().IntVarP(&lines, "tail", "t", 0, "Number of lines from the end of the log file to display")
	daemonLogsCmd.Flags().BoolVarP(&clear, "clear", "c", false, "Clear the log file")
	daemonLogsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow the log file (like tail -f)")
}

func showDaemonLogs() {
	if follow {
		followDaemonLogs()
		return
	}

	logData, err := os.ReadFile(logFilePath)
	if err != nil {
		logger.LogError("Failed to read log file: " + err.Error())
		return
	}

	if lines > 0 {
		logLines := strings.Split(string(logData), "\n")
		if len(logLines) > lines {
			logLines = logLines[len(logLines)-lines:]
		}
		logger.LogInfo(strings.Join(logLines, "\n"))
	} else {
		logger.LogInfo(string(logData))
	}
}

func clearDaemonLogs() {
	err := os.Truncate(logFilePath, 0)
	if err != nil {
		logger.LogError("Failed to clear log file: " + err.Error())
		return
	}
	logger.LogInfo("Log file cleared.")
}

func followDaemonLogs() {
	file, err := os.Open(logFilePath)
	if err != nil {
		logger.LogError("Failed to open log file: " + err.Error())
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	stat, err := file.Stat()
	if err != nil {
		logger.LogError("Failed to get file stat: " + err.Error())
		return
	}

	start := stat.Size()
	_, err = file.Seek(start, 0)
	if err != nil {
		logger.LogError("Failed to seek file: " + err.Error())
		return
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Print(line)
	}
}
