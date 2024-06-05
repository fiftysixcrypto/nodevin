package bitcoin

import (
	"github.com/curveballdaniel/nodevin/pkg/docker/compose"
	"github.com/spf13/viper"
)

func CreateBitcoinComposeFile(cwd string) (string, error) {
	bitcoinNetwork := viper.GetString("network")
	baseTestnet := viper.GetBool("testnet")

	var network string

	if bitcoinNetwork != "" {
		network = "bitcoin-" + bitcoinNetwork
	} else if baseTestnet {
		network = "bitcoin-testnet"
	} else {
		network = "bitcoin"
	}

	bitcoinBaseComposeConfig, err := compose.GetNetworkComposeConfig(network)
	if err != nil {
		return "", err
	}

	composeFilePath, err := compose.CreateComposeFile("bitcoin-core", cwd, bitcoinBaseComposeConfig)
	if err != nil {
		return "", err
	}

	return composeFilePath, nil
}
