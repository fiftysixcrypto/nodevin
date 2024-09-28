/*
// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2024 The Nodevin Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
*/

package root

import (
	"fmt"

	"github.com/fiftysixcrypto/nodevin/internal/version"
	"github.com/fiftysixcrypto/nodevin/pkg/blockchain"
	"github.com/fiftysixcrypto/nodevin/pkg/daemon"
	"github.com/fiftysixcrypto/nodevin/pkg/initialize"
	"github.com/fiftysixcrypto/nodevin/pkg/update"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nodevin",
	Short: "nodevin CLI",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Nodevin",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("nodevin CLI v" + version.Version)
	},
}

func init() {
	// Define flags and configuration settings
	rootCmd.PersistentFlags().String("command", "", "Initial command to run the node")
	rootCmd.PersistentFlags().Bool("testnet", false, "Run assumed network testnet")
	rootCmd.PersistentFlags().String("network", "", "Run node attached to a specific network (network name -- ex: goerli, testnet3)")
	rootCmd.PersistentFlags().String("rpc-user", "user", "Username passed in via command for JSON RPC -- (default: user)")
	rootCmd.PersistentFlags().String("rpc-pass", "fiftysix", "Username passed in via command for JSON RPC -- (default: fiftysix)")
	rootCmd.PersistentFlags().Bool("cookie-auth", false, "Use authentication directly with node cookie file -- (default: false)")
	rootCmd.PersistentFlags().StringSlice("ports", []string{}, "Ports to bind to the node")
	rootCmd.PersistentFlags().StringSlice("volumes", []string{}, "Docker volumes to mount for compose file")
	rootCmd.PersistentFlags().StringSlice("volume-definitions", []string{}, "Docker volume definitions for compose file")
	rootCmd.PersistentFlags().String("image", "", "Docker image to use for the node (image name -- ex: fiftysix/bitcoin-core)")
	rootCmd.PersistentFlags().String("version", "", "Version of Docker image to use for the node (tag -- ex: latest, 27.0)")
	rootCmd.PersistentFlags().String("restart", "no", "Whether or not to restart software on failure (docker restart parameter -- ex: always, no)")
	rootCmd.PersistentFlags().String("container-name", "", "Docker container name for compose file")
	rootCmd.PersistentFlags().StringSlice("docker-networks", []string{}, "Docker networks to connect to for compose file")
	rootCmd.PersistentFlags().String("network-driver", "", "Docker network driver for compose file")
	rootCmd.PersistentFlags().StringToString("volume-labels", map[string]string{}, "Docker volume labels for compose file")

	rootCmd.PersistentFlags().String("cpu-limit", "", "Maximum CPU limit of use (amount of CPUs -- ex: 1.5)")
	rootCmd.PersistentFlags().String("mem-limit", "", "Maximum memory limit of use (positive integer followed by 'b', 'k', 'm', 'g', to indicate bytes, kilobytes, megabytes, or gigabytes -- ex: 50m)")
	rootCmd.PersistentFlags().String("cpu-reservation", "", "Reserve a set amount of CPU for use (amount of CPUs -- ex: 1.5)")
	rootCmd.PersistentFlags().String("mem-reservation", "", "Reserve a set amount of memory for use (positive integer followed by 'b', 'k', 'm', 'g', to indicate bytes, kilobytes, megabytes, or gigabytes -- ex: 50m)")

	// Bitcoin specific flags
	rootCmd.PersistentFlags().Bool("ord", false, "Run ordinal software ord alongside the Bitcoin/Litecoin node")
	rootCmd.PersistentFlags().String("ord-image", "fiftysix/ord", "Docker image to use for ord (image name -- ex: fiftysix/ord)")
	rootCmd.PersistentFlags().String("ord-version", "latest", "Version of Docker image to use for ord (tag -- ex: latest, 27.0)")

	// Litecoin specific flags
	rootCmd.PersistentFlags().Bool("ord-litecoin", false, "Run ordinal software ord alongside the Litecoin node")
	rootCmd.PersistentFlags().String("ord-litecoin-image", "fiftysix/ord-litecoin", "Docker image to use for ord (image name -- ex: fiftysix/ord-litecoin)")
	rootCmd.PersistentFlags().String("ord-litecoin-version", "latest", "Version of Docker image to use for ord (tag -- ex: latest, 27.0)")

	// Bind flags to viper
	viper.BindPFlag("command", rootCmd.PersistentFlags().Lookup("command"))
	viper.BindPFlag("testnet", rootCmd.PersistentFlags().Lookup("testnet"))
	viper.BindPFlag("network", rootCmd.PersistentFlags().Lookup("network"))
	viper.BindPFlag("rpc-user", rootCmd.PersistentFlags().Lookup("rpc-user"))
	viper.BindPFlag("rpc-pass", rootCmd.PersistentFlags().Lookup("rpc-pass"))
	viper.BindPFlag("cookie-auth", rootCmd.PersistentFlags().Lookup("cookie-auth"))
	viper.BindPFlag("ports", rootCmd.PersistentFlags().Lookup("ports"))
	viper.BindPFlag("volumes", rootCmd.PersistentFlags().Lookup("volumes"))
	viper.BindPFlag("volume-definitions", rootCmd.PersistentFlags().Lookup("volume-definitions"))
	viper.BindPFlag("image", rootCmd.PersistentFlags().Lookup("image"))
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
	viper.BindPFlag("restart", rootCmd.PersistentFlags().Lookup("restart"))
	viper.BindPFlag("container-name", rootCmd.PersistentFlags().Lookup("container-name"))
	viper.BindPFlag("docker-networks", rootCmd.PersistentFlags().Lookup("docker-networks"))
	viper.BindPFlag("network-driver", rootCmd.PersistentFlags().Lookup("network-driver"))
	viper.BindPFlag("volume-labels", rootCmd.PersistentFlags().Lookup("volume-labels"))

	viper.BindPFlag("cpu-limit", rootCmd.PersistentFlags().Lookup("cpu-limit"))
	viper.BindPFlag("mem-limit", rootCmd.PersistentFlags().Lookup("mem-limit"))
	viper.BindPFlag("cpu-reservation", rootCmd.PersistentFlags().Lookup("cpu-reservation"))
	viper.BindPFlag("mem-reservation", rootCmd.PersistentFlags().Lookup("mem-reservation"))

	// Bitcoin specific flags
	viper.BindPFlag("ord", rootCmd.PersistentFlags().Lookup("ord"))
	viper.BindPFlag("ord-image", rootCmd.PersistentFlags().Lookup("ord-image"))
	viper.BindPFlag("ord-version", rootCmd.PersistentFlags().Lookup("ord-version"))

	// Litecoin specific flags
	viper.BindPFlag("ord-litecoin", rootCmd.PersistentFlags().Lookup("ord-litecoin"))
	viper.BindPFlag("ord-litecoin-image", rootCmd.PersistentFlags().Lookup("ord-litecoin-image"))
	viper.BindPFlag("ord-litecoin-version", rootCmd.PersistentFlags().Lookup("ord-litecoin-version"))

	// Add blockchain commands
	rootCmd.AddCommand(blockchain.RequestCmd)
	rootCmd.AddCommand(blockchain.BackupCmd)
	rootCmd.AddCommand(blockchain.RestartNodeCmd)
	rootCmd.AddCommand(blockchain.ShellCmd)
	rootCmd.AddCommand(blockchain.DeleteVolumeCmd)
	rootCmd.AddCommand(blockchain.StartNodeCmd)
	rootCmd.AddCommand(blockchain.StopNodeCmd)
	rootCmd.AddCommand(blockchain.LogsCmd)
	rootCmd.AddCommand(blockchain.InfoCmd)

	// Add init command
	rootCmd.AddCommand(initialize.InitCmd)

	// Add manual update commands
	rootCmd.AddCommand(update.UpdateCmd)

	// Add daemon commands
	rootCmd.AddCommand(daemon.DaemonCmd)
}

func Execute() error {
	rootCmd.AddCommand(versionCmd)
	return rootCmd.Execute()
}
