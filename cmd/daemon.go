package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the NodeVin daemon",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		runDaemon()
	},
}

func runDaemon() {
	logInfo("Starting NodeVin daemon...")

	// Set up a signal channel to capture OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Set up a ticker to check for updates periodically
	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for range ticker.C {
			logInfo("Checking for updates...")
			if err := checkForUpdates(); err != nil {
				logError("Failed to check for updates: " + err.Error())
				continue
			}
			logInfo("Update downloaded. Applying update...")
			if err := applyUpdate(); err != nil {
				logError("Failed to apply update: " + err.Error())
				continue
			}
			logInfo("Update applied successfully. Please restart the application.")
		}
	}()

	// Block until a signal is received
	sig := <-sigs
	logInfo("Received signal: " + sig.String())

	logInfo("Shutting down NodeVin daemon...")
}

func init() {
	daemonCmd.Flags().StringVar(&config.Network, "network", "mainnet", "Blockchain network to connect to")
	daemonCmd.Flags().StringVar(&config.StoragePath, "storage_path", "/var/lib/nodevin", "Path to store blockchain data")
	daemonCmd.Flags().IntVar(&config.Port, "port", 30303, "Port to bind the node")
	daemonCmd.Flags().StringVar(&config.ResourceLimit, "resource_limit", "2GB", "Resource limit for the node")

	viper.BindPFlag("network", daemonCmd.Flags().Lookup("network"))
	viper.BindPFlag("storage_path", daemonCmd.Flags().Lookup("storage_path"))
	viper.BindPFlag("port", daemonCmd.Flags().Lookup("port"))
	viper.BindPFlag("resource_limit", daemonCmd.Flags().Lookup("resource_limit"))

	rootCmd.AddCommand(daemonCmd)
}
