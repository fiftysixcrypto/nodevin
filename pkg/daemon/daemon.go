package daemon

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the NodeVin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		runDaemon()
	},
}

func runDaemon() {
	logger.LogInfo("Starting NodeVin daemon...")

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Set up a ticker to check for updates periodically
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			logger.LogInfo("Checking for updates...")
			if err := checkForUpdates(); err != nil {
				logger.LogError("Failed to check for updates: " + err.Error())
				continue
			}
			logger.LogInfo("Update downloaded. Applying update...")
			if err := applyUpdate(); err != nil {
				logger.LogError("Failed to apply update: " + err.Error())
				continue
			}
			logger.LogInfo("Update applied successfully. Please restart the application.")
		}
	}()

	// Block until a signal is received
	sig := <-sigs
	logger.LogInfo("Received signal: " + sig.String())

	// Perform cleanup and shutdown
	shutdownDaemon()
	logger.LogInfo("Shutting down NodeVin daemon...")
}

func shutdownDaemon() {
	// Implement the cleanup and shutdown logic here
}

func Execute() error {
	return Cmd.Execute()
}
