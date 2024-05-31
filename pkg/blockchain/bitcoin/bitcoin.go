package bitcoin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func CreateBitcoinEnv() (string, error) {
	port := viper.GetString("port")
	storagePath := viper.GetString("data-dir")
	extraArgs := viper.GetString("args")

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Specify the path to the custom .env file in the current working directory
	envFilePath := filepath.Join(cwd, ".bitcoin.env")

	// Create or open the custom .env file
	file, err := os.Create(envFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create .bitcoin.env file: %w", err)
	}
	defer file.Close()

	// Write environment variables to the custom .env file
	if len(port) > 0 {
		_, err = fmt.Fprintf(file, "FIFTYSIX_BITCOIN_RPC_PORT=%s\n", port)
		if err != nil {
			return "", fmt.Errorf("failed to write to .bitcoin.env file: %w", err)
		}
	}

	if len(storagePath) > 0 {
		_, err = fmt.Fprintf(file, "FIFTYSIX_BITCOIN_STORAGE_PATH=%s\n", storagePath)
		if err != nil {
			return "", fmt.Errorf("failed to write to .bitcoin.env file: %w", err)
		}
	}

	if len(extraArgs) > 0 {
		_, err = fmt.Fprintf(file, "FIFTYSIX_BITCOIN_EXTRA_ARGS=%s\n", extraArgs)
		if err != nil {
			return "", fmt.Errorf("failed to write to .bitcoin.env file: %w", err)
		}
	}

	/*
	     cpus: ${FIFTYSIX_BITCOIN_LIMITS_CPUS}
	     memory: ${FIFTYSIX_BITCOIN_LIMITS_MEMORY}
	   reservations:
	     cpus: ${FIFTYSIX_BITCOIN_RESERVATIONS_CPUS}
	     memory: ${FIFTYSIX_BITCOIN_RESERVATIONS_MEMORY}
	*/

	return envFilePath, nil
}
