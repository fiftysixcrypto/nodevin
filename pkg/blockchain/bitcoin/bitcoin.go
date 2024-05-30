package bitcoin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/config"
	"github.com/curveballdaniel/nodevin/internal/logger"
	"github.com/curveballdaniel/nodevin/pkg/blockchain"
)

type Bitcoin struct{}

func (b Bitcoin) StartNode(config config.Config) error {
	fmt.Println("Starting Bitcoin node with config:", config)

	// Set environment variables for Docker Compose
	//os.Setenv("BITCOIN_PORT", fmt.Sprintf("%d", config.Port))
	//os.Setenv("BITCOIN_STORAGE_PATH", config.StoragePath)
	//os.Setenv("BITCOIN_EXTRA_ARGS", fmt.Sprintf("%s", extraArgs))

	//     command: bitcoind ${BITCOIN_EXTRA_ARGS}
	//  - "${FIFTYSIX_BITCOIN_RPC_PORT:-8332}:${FIFTYSIX_BITCOIN_RPC_PORT:-8332}"

	// Define the path to the Docker Compose file
	composeFilePath := "docker/bitcoin/docker-compose_bitcoin-core.yml"

	// Start the services using Docker Compose
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "up", "-d")
	cmd.Stdout = logWriter{}
	cmd.Stderr = logWriter{}

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to start Docker Compose services: " + err.Error())
		return err
	}

	logger.LogInfo("Bitcoin node started successfully")

	return nil
}

func (b Bitcoin) StopNode() error {
	fmt.Println("Stopping Bitcoin node")

	// Define the path to the Docker Compose file
	composeFilePath := "docker/bitcoin/docker-compose_bitcoin-core.yml"

	// Stop the services using Docker Compose
	cmd := exec.Command("docker-compose", "-f", composeFilePath, "down")
	cmd.Stdout = logWriter{}
	cmd.Stderr = logWriter{}

	if err := cmd.Run(); err != nil {
		logger.LogError("Failed to stop Docker Compose services: " + err.Error())
		return err
	}

	logger.LogInfo("Bitcoin node stopped successfully")

	return nil
}

func init() {
	blockchain.RegisterBlockchain("bitcoin", Bitcoin{})
}

type logWriter struct{}

func (f logWriter) Write(bytes []byte) (int, error) {
	return fmt.Fprint(os.Stdout, string(bytes))
}
