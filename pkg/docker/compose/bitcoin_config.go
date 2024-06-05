package compose

import (
	"fmt"
)

func GetBitcoinNetworkComposeConfig(network string) (NetworkConfig, error) {
	baseConfig := NetworkConfig{
		Image:    "fiftysix/bitcoin-core",
		Version:  "latest",
		Ports:    []string{"8332:8332", "8333:8333"},
		Volumes:  []string{"bitcoin-core-data:/node/bitcoin-core"},
		Networks: []string{"bitcoin-net"},
		NetworkDefs: map[string]NetworkDetails{
			"bitcoin-net": {
				Driver: "bridge",
			},
		},
		VolumeDefs: map[string]VolumeDetails{
			"bitcoin-core-data": {
				Labels: map[string]string{
					"blockchain.software": "bitcoin-core",
				},
			},
		},
	}

	switch network {
	case "bitcoin":
		baseConfig.ContainerName = "bitcoin-core"
		baseConfig.Command = "bitcoind --server=1 --rpcbind=0.0.0.0 --rpcport=8332 --rpcallowip=0.0.0.0/0"
	case "bitcoin-testnet":
		baseConfig.ContainerName = "bitcoin-core-testnet"
		baseConfig.Command = "bitcoind --testnet --server=1 --rpcbind=0.0.0.0 --rpcport=18332 --rpcallowip=0.0.0.0/0"
		baseConfig.Ports = []string{"18332:18332", "18333:18333"}
		baseConfig.Networks = []string{"bitcoin-testnet-net"}
		baseConfig.NetworkDefs = map[string]NetworkDetails{
			"bitcoin-testnet-net": {
				Driver: "bridge",
			},
		}
		baseConfig.Volumes = []string{"bitcoin-core-testnet-data:/node/bitcoin-core"}
		baseConfig.VolumeDefs = map[string]VolumeDetails{
			"bitcoin-core-testnet-data": {
				Labels: map[string]string{
					"blockchain.software": "bitcoin-core",
				},
			},
		}
	default:
		return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
	}

	return baseConfig, nil
}
