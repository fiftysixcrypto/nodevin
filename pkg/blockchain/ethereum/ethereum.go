package bitcoin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fiftysixcrypto/nodevin/internal/config"
	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/pkg/blockchain"
)

type Bitcoin struct{}

func (b Bitcoin) StartNode(config config.Config) error {
	fmt.Println("Starting Bitcoin node with config:", config)

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
