package root

import (
	"fmt"

	"github.com/curveballdaniel/nodevin/internal/version"
	"github.com/curveballdaniel/nodevin/pkg/blockchain"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nodevin",
	Short: "Nodevin CLI",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of NodeVin",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nodevin CLI v" + version.Version)
	},
}

func init() {
	// Define flags and configuration settings
	rootCmd.PersistentFlags().String("command", "", "Initial command to run the node")
	rootCmd.PersistentFlags().StringSlice("ports", []string{}, "Ports to bind to the node")
	rootCmd.PersistentFlags().StringSlice("volumes", []string{}, "Docker volumes to mount for compose file")
	rootCmd.PersistentFlags().StringSlice("volume-definitions", []string{}, "Docker volume definitions for compose file")
	rootCmd.PersistentFlags().String("image", "", "Docker image to use for the node (image name -- ex: fiftysix/bitcoin-core)")
	rootCmd.PersistentFlags().String("version", "", "Version of Docker image to use for the node (tag -- ex: latest, 27.0)")
	rootCmd.PersistentFlags().String("container-name", "", "Docker container name for compose file")
	rootCmd.PersistentFlags().StringSlice("networks", []string{}, "Docker networks to connect to for compose file")
	rootCmd.PersistentFlags().String("network-driver", "", "Docker network driver for compose file")
	rootCmd.PersistentFlags().StringToString("volume-labels", map[string]string{}, "Docker volume labels for compose file")

	rootCmd.PersistentFlags().String("cpu-limit", "", "Maximum CPU limit of use (amount of CPUs -- ex: 1.5)")
	rootCmd.PersistentFlags().String("mem-limit", "", "Maximum memory limit of use (positive integer followed by 'b', 'k', 'm', 'g', to indicate bytes, kilobytes, megabytes, or gigabytes -- ex: 50m)")
	rootCmd.PersistentFlags().String("cpu-reservation", "", "Reserve a set amount of CPU for use (amount of CPUs -- ex: 1.5)")
	rootCmd.PersistentFlags().String("mem-reservation", "", "Reserve a set amount of memory for use (positive integer followed by 'b', 'k', 'm', 'g', to indicate bytes, kilobytes, megabytes, or gigabytes -- ex: 50m)")

	// Bind flags to viper
	viper.BindPFlag("command", rootCmd.PersistentFlags().Lookup("command"))
	viper.BindPFlag("ports", rootCmd.PersistentFlags().Lookup("ports"))
	viper.BindPFlag("volumes", rootCmd.PersistentFlags().Lookup("volumes"))
	viper.BindPFlag("volume-definitions", rootCmd.PersistentFlags().Lookup("volume-definitions"))
	viper.BindPFlag("image", rootCmd.PersistentFlags().Lookup("image"))
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
	viper.BindPFlag("container-name", rootCmd.PersistentFlags().Lookup("container-name"))
	viper.BindPFlag("networks", rootCmd.PersistentFlags().Lookup("networks"))
	viper.BindPFlag("network-driver", rootCmd.PersistentFlags().Lookup("network-driver"))
	viper.BindPFlag("volume-labels", rootCmd.PersistentFlags().Lookup("volume-labels"))

	viper.BindPFlag("cpu-limit", rootCmd.PersistentFlags().Lookup("cpu-limit"))
	viper.BindPFlag("mem-limit", rootCmd.PersistentFlags().Lookup("mem-limit"))
	viper.BindPFlag("cpu-reservation", rootCmd.PersistentFlags().Lookup("cpu-reservation"))
	viper.BindPFlag("mem-reservation", rootCmd.PersistentFlags().Lookup("mem-reservation"))

	// Add blockchain commands
	rootCmd.AddCommand(blockchain.ShellCmd)
	rootCmd.AddCommand(blockchain.DeleteVolumeCmd)
	rootCmd.AddCommand(blockchain.StartNodeCmd)
	rootCmd.AddCommand(blockchain.StopNodeCmd)
	rootCmd.AddCommand(blockchain.LogsCmd)
	rootCmd.AddCommand(blockchain.InfoCmd)
}

func Execute() error {
	rootCmd.AddCommand(versionCmd)
	return rootCmd.Execute()
}
