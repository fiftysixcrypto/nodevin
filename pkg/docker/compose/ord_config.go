package compose

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetOrdNetworkComposeConfig(network string) (NetworkConfig, error) {
	baseConfig := NetworkConfig{
		Image:    "fiftysix/ord",
		Version:  "latest",
		Ports:    []string{"80:80"},
		Volumes:  []string{"bitcoin-core-data:/node/bitcoin-core", "ord-data:/node/ord"},
		Networks: []string{"bitcoin-net"},
		NetworkDefs: map[string]NetworkDetails{
			"bitcoin-net": {
				Driver: "bridge",
			},
		},
		VolumeDefs: map[string]VolumeDetails{
			"ord-data": {
				Labels: map[string]string{
					"nodevin.blockchain.software": "ord",
				},
			},
		},
	}

	switch network {
	case "ord":
		baseConfig.ContainerName = "ord"
		baseConfig.Command = "ord --bitcoin-rpc-url http://bitcoin-core:8332"
	case "ord-testnet":
		baseConfig.ContainerName = "ord-testnet"
		baseConfig.Command = "ord --testnet --bitcoin-rpc-url http://bitcoin-core:18332"
		baseConfig.Volumes = []string{"bitcoin-core-testnet-data:/node/bitcoin-core", "ord-testnet-data:/node/ord"}
		baseConfig.Networks = []string{"bitcoin-testnet-net"}
		baseConfig.NetworkDefs = map[string]NetworkDetails{
			"bitcoin-testnet-net": {
				Driver: "bridge",
			},
		}
		baseConfig.VolumeDefs = map[string]VolumeDetails{
			"ord-testnet-data": {
				Labels: map[string]string{
					"nodevin.blockchain.software": "ord",
				},
			},
		}
	default:
		return NetworkConfig{}, fmt.Errorf("unknown network: %s", network)
	}

	cookieAuth := viper.GetBool("ord-cookie-auth")

	if !cookieAuth {
		// Add RPC user/pass to command
		rpcUsername := viper.GetString("ord-rpc-user")
		rpcPassword := viper.GetString("ord-rpc-pass")

		if rpcUsername == "" {
			rpcUsername = "user"
		}

		if rpcPassword == "" {
			rpcPassword = "fiftysix"
		}

		baseConfig.Command = baseConfig.Command + " " + fmt.Sprintf("--bitcoin-rpc-username %s", rpcUsername) + " " + fmt.Sprintf("--bitcoin-rpc-password %s", rpcPassword)
	}

	baseConfig.Command = baseConfig.Command + " server"

	return baseConfig, nil
}
