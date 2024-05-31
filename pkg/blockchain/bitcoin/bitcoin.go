package bitcoin

import (
	"github.com/curveballdaniel/nodevin/pkg/docker/compose"
)

func CreateBitcoinComposeFile(cwd string) (string, error) {
	bitcoinBaseComposeConfig, err := compose.GetNetworkComposeConfig("bitcoin")
	if err != nil {
		return "", err
	}

	composeFilePath, err := compose.CreateComposeFile("bitcoin-core", cwd, bitcoinBaseComposeConfig)
	if err != nil {
		return "", err
	}

	return composeFilePath, nil
}
