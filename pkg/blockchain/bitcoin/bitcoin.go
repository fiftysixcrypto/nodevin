package bitcoin

import (
	"github.com/curveballdaniel/nodevin/internal/utils"
	"github.com/curveballdaniel/nodevin/pkg/docker/compose"
	"github.com/spf13/viper"
)

func CreateBitcoinComposeFile(cwd string) (string, error) {
	var network string

	if utils.CheckIfTestnetOrTestnetNetworkFlag() {
		network = "bitcoin-testnet"
	} else {
		network = "bitcoin"
	}

	bitcoinBaseComposeConfig, err := compose.GetBitcoinNetworkComposeConfig(network)
	if err != nil {
		return "", err
	}

	if viper.GetBool("ord") {
		var ordNetwork string

		if utils.CheckIfTestnetOrTestnetNetworkFlag() {
			ordNetwork = "ord-testnet"
		} else {
			ordNetwork = "ord"
		}

		ordComposeConfig, err := compose.GetOrdNetworkComposeConfig(ordNetwork)
		if err != nil {
			return "", err
		}

		composeFilePath, err := compose.CreateComposeFile(
			bitcoinBaseComposeConfig.ContainerName,
			bitcoinBaseComposeConfig,
			[]string{"ord"},
			[]compose.NetworkConfig{ordComposeConfig},
			cwd)

		if err != nil {
			return "", err
		}

		return composeFilePath, nil
	} else {
		composeFilePath, err := compose.CreateComposeFile(
			bitcoinBaseComposeConfig.ContainerName,
			bitcoinBaseComposeConfig,
			[]string{},
			[]compose.NetworkConfig{},
			cwd)

		if err != nil {
			return "", err
		}

		return composeFilePath, nil
	}
}
