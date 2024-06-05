package compose

import (
	"fmt"
)

func GetNetworkComposeConfig(network string) (NetworkConfig, error) {
	switch {
	case network == "bitcoin" || network == "bitcoin-testnet":
		return GetBitcoinNetworkComposeConfig(network)
	case network == "ethereum":
		return GetEthereumNetworkComposeConfig(network)
	default:
		return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
	}
}
