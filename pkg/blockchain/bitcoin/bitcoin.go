package bitcoin

import (
	"github.com/curveballdaniel/nodevin/internal/utils"
	"github.com/curveballdaniel/nodevin/pkg/docker/compose"
)

func CreateBitcoinComposeFile(cwd string) (string, error) {
	var network string

	if utils.CheckIfTestnetOrTestnetNetwork() {
		network = "bitcoin-testnet"
	} else {
		network = "bitcoin"
	}

	bitcoinBaseComposeConfig, err := compose.GetNetworkComposeConfig(network)
	if err != nil {
		return "", err
	}

	composeFilePath, err := compose.CreateComposeFile(bitcoinBaseComposeConfig.ContainerName, cwd, bitcoinBaseComposeConfig)
	if err != nil {
		return "", err
	}

	return composeFilePath, nil
}
