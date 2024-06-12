package daemon

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/update"
	"github.com/spf13/cobra"
)

var DaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the nodevin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		runDaemon()
	},
}

func runDaemon() {
	logger.LogInfo("Starting nodevin daemon...")

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Set up a ticker to check for updates periodically
	updateTicker := time.NewTicker(24 * time.Hour)
	imageCheckTicker := time.NewTicker(1 * time.Minute) //time.Hour

	go func() {
		for {
			select {
			case <-updateTicker.C:
				update.CommandCheckForUpdatesWorkflow()
			case <-imageCheckTicker.C:
				update.CommandCheckAndUpdateDockerImagesWorkflow()
			}
		}
	}()

	logger.LogInfo("Successfully started daemon. Nodevin will monitor your running nodes and update them if a newer version of the software is launched.")

	// Block until a signal is received
	sig := <-sigs
	logger.LogInfo("Received signal: " + sig.String())

	// Perform cleanup and shutdown
	shutdownDaemon()
	logger.LogInfo("Shutting down nodevin daemon...")
}

func shutdownDaemon() {
	// Implement the cleanup and shutdown logic here
}
