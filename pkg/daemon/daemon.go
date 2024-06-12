package daemon

import (
	"github.com/spf13/cobra"
)

var DaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage the nodevin daemon",
}

func init() {
	DaemonCmd.AddCommand(daemonStartCmd, daemonStopCmd, daemonLogsCmd)
}
