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
	rootCmd.PersistentFlags().String("port", "", "Port to bind the node")
	rootCmd.PersistentFlags().String("data-dir", "", "Path to store blockchain data")
	rootCmd.PersistentFlags().String("args", "", "Extra flags/arguments to add to node execution")

	rootCmd.PersistentFlags().String("cpu-limit", "", "Maximum CPU limit of use (amount of CPUs -- ex 1.5)")
	rootCmd.PersistentFlags().String("mem-limit", "", "Maximum memory limit of use (positive integer followed by 'b', 'k', 'm', 'g', to indicate bytes, kilobytes, megabytes, or gigabytes -- ex 50m)")
	rootCmd.PersistentFlags().String("cpu-reservation", "", "Reserve a set amount of CPU for use (amount of CPUs -- ex 1.5)")
	rootCmd.PersistentFlags().String("mem-reservation", "", "Reserve a set amount of memory for use (positive integer followed by 'b', 'k', 'm', 'g', to indicate bytes, kilobytes, megabytes, or gigabytes -- ex 50m)")

	// Bind flags to viper
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("data-dir", rootCmd.PersistentFlags().Lookup("data-dir"))
	viper.BindPFlag("args", rootCmd.PersistentFlags().Lookup("args"))

	viper.BindPFlag("cpu-limit", rootCmd.PersistentFlags().Lookup("cpu-limit"))
	viper.BindPFlag("mem-limit", rootCmd.PersistentFlags().Lookup("mem-limit"))
	viper.BindPFlag("cpu-reservation", rootCmd.PersistentFlags().Lookup("cpu-reservation"))
	viper.BindPFlag("mem-reservation", rootCmd.PersistentFlags().Lookup("mem-reservation"))

	// Add blockchain commands
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
