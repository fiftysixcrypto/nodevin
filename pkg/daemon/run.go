package daemon

import (
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/update"
)

func runDaemon() {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogError("Failed to open log file: " + err.Error())
		return
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(multiWriter)

	logger.LogInfo("")
	logger.LogInfo("Starting nodevin daemon...")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	updateTicker := time.NewTicker(24 * time.Hour)
	imageCheckTicker := time.NewTicker(1 * time.Minute)

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

	logger.LogInfo("Successfully started daemon.")
	logger.LogInfo("Nodevin monitors running nodes and updates them if newer versions are released.")

	sig := <-sigs
	logger.LogInfo("Received signal: " + sig.String())

	logger.LogInfo("Shutting down nodevin daemon...")
}
