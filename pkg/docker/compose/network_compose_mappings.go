package compose

import "fmt"

func GetNetworkComposeConfig(network string) (NetworkConfig, error) {
	switch network {
	case "bitcoin":
		return NetworkConfig{
			Image:         "fiftysix/bitcoin-core:latest",
			ContainerName: "bitcoin-core",
			Command:       "bitcoind",
			Ports:         []string{"8332:8332", "8333:8333"},
			Volumes:       []string{"bitcoin-core-data:/node/bitcoin-core"},
			Networks:      []string{"bitcoin-net"},
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
		}, nil
	case "ethereum":
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
	default:
		return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
	}
}
