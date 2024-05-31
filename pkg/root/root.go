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

	// Bind flags to viper
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("data-dir", rootCmd.PersistentFlags().Lookup("data-dir"))
	viper.BindPFlag("args", rootCmd.PersistentFlags().Lookup("args"))

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
