package daemon

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

const pidFilePath = "/tmp/nodevin_daemon.pid"
const logFilePath = "/tmp/nodevin_daemon.log"

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the nodevin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		detached, _ := cmd.Flags().GetBool("detach")
		if detached {
			startDaemonDetached()
		} else {
			runDaemon()
		}
	},
}

func init() {
	daemonStartCmd.Flags().BoolP("detach", "d", false, "Run daemon in the background")
}

func startDaemonDetached() {
	if _, err := os.Stat(pidFilePath); err == nil {
		logger.LogError("Daemon is already running.")
		return
	}

	cmd := exec.Command(os.Args[0], "daemon", "start")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Start()

	err := os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
	if err != nil {
		logger.LogError("Failed to write PID file: " + err.Error())
		return
	}

	logger.LogInfo(fmt.Sprintf("Daemon started with PID: %v", cmd.Process.Pid))
	fmt.Println("Check daemon logs: nodevin daemon logs")
	fmt.Println("Stop daemon: nodevin daemon stop")
}
