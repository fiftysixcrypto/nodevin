package compose

import (
	"fmt"
)

func GetEthereumNetworkComposeConfig(network string) (NetworkConfig, error) {
	if network == "ethereum" {
		return NetworkConfig{
			Image:         "ethereum/client-go:latest",
			ContainerName: "ethereum-node",
			Command:       "geth",
			Ports:         []string{"8545:8545", "30303:30303"},
			Volumes:       []string{"ethereum-node-data:/root/.ethereum"},
			Networks:      []string{"ethereum-net"},
			NetworkDefs: map[string]NetworkDetails{
				"ethereum-net": {
					Driver: "bridge",
				},
			},
			VolumeDefs: map[string]VolumeDetails{
				"ethereum-node-data": {
					Labels: map[string]string{
						"blockchain.software": "ethereum-node",
					},
				},
			},
		}, nil
	}
	return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
}
