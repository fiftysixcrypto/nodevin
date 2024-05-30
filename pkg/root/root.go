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
	Short: "NodeVin CLI",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of NodeVin",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NodeVin CLI v" + version.Version)
	},
}

func init() {
	// Define flags and configuration settings
	rootCmd.PersistentFlags().String("network", "bitcoin", "Blockchain network to connect to")
	rootCmd.PersistentFlags().String("storage_path", "/var/lib/nodevin", "Path to store blockchain data")
	rootCmd.PersistentFlags().Int("port", 30303, "Port to bind the node")
	rootCmd.PersistentFlags().String("resource_limit", "2GB", "Resource limit for the node")

	// Bind flags to viper
	viper.BindPFlag("network", rootCmd.PersistentFlags().Lookup("network"))
	viper.BindPFlag("storage_path", rootCmd.PersistentFlags().Lookup("storage_path"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("resource_limit", rootCmd.PersistentFlags().Lookup("resource_limit"))

	// Add blockchain commands
	rootCmd.AddCommand(blockchain.StartNodeCmd)
	rootCmd.AddCommand(blockchain.StopNodeCmd)
	rootCmd.AddCommand(blockchain.LogsCmd)
	rootCmd.AddCommand(blockchain.InfoCmd)
}

func Execute() error {
	rootCmd.AddCommand(versionCmd)
	return rootCmd.Execute()
}
