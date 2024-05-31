package bitcoin

import (
	"github.com/curveballdaniel/nodevin/pkg/docker/compose"
	"github.com/spf13/viper"
)

func CreateBitcoinComposeFile(cwd string) (string, error) {
	/*
			viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
			viper.BindPFlag("data-dir", rootCmd.PersistentFlags().Lookup("data-dir"))
			viper.BindPFlag("args", rootCmd.PersistentFlags().Lookup("args"))

		port := viper.GetString("port")
	*/
	// Define Viper configuration keys and set defaults where necessary

	viper.SetDefault("command", "bitcoind")
	viper.SetDefault("ports", []string{"8332:8332", "8333:8333"})

	viper.SetDefault("volumes", []string{"bitcoin-core-data:/node/bitcoin-core"})
	viper.SetDefault("volume-definitions", []string{"bitcoin-core-data"})

	viper.SetDefault("image", "fiftysix/bitcoin-core:latest")

	viper.SetDefault("container_name", "bitcoin-core")
	viper.SetDefault("networks", []string{"bitcoin-net"})
	viper.SetDefault("network-driver", "bridge")
	viper.SetDefault("volume-labels", map[string]string{
		"blockchain.software": "bitcoin-core",
	})

	// Define the override configuration using Viper
	override := compose.NetworkConfig{
		Image:         viper.GetString("image"),
		ContainerName: viper.GetString("container_name"),
		Command:       viper.GetString("command"),
		Ports:         viper.GetStringSlice("ports"),
		Volumes:       viper.GetStringSlice("volumes"),
		Networks:      viper.GetStringSlice("networks"),
		Deploy: compose.Deploy{
			Resources: compose.Resources{
				Limits: compose.ResourceDetails{
					CPUs:   viper.GetString("cpu-limit"),
					Memory: viper.GetString("mem-limit"),
				},
				Reservations: compose.ResourceDetails{
					CPUs:   viper.GetString("cpu-reservation"),
					Memory: viper.GetString("mem-reservation"),
				},
			},
		},
		NetworkDefs: map[string]compose.NetworkDetails{
			viper.GetStringSlice("networks")[0]: {
				Driver: viper.GetString("network-driver"),
			},
		},
		VolumeDefs: map[string]compose.VolumeDetails{
			viper.GetStringSlice("volume-definitions")[0]: {
				Labels: viper.GetStringMapString("volume-labels"),
			},
		},
	}

	bitcoinBaseComposeConfig, err := compose.GetNetworkComposeConfig("bitcoin")
	if err != nil {
		return "", err
	}

	composeFilePath, err := compose.CreateComposeFile("bitcoin-core", cwd, bitcoinBaseComposeConfig, override)
	if err != nil {
		return "", err
	}

	return composeFilePath, nil
}
