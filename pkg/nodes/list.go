package nodes

import (
	"fmt"
	"sort"

	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported networks",
	Run: func(cmd *cobra.Command, args []string) {
		listAllNetworks()
	},
}

func listAllNetworks() {
	var networkNames []string
	for network := range utils.NetworkContainerMap() {
		networkNames = append(networkNames, network)
	}

	sort.Strings(networkNames)

	fmt.Printf("Supported networks: %s\n", utils.GetCommandSupportedNetworks())
	fmt.Print("\nHelpful Commands:\n")
	fmt.Printf("%s start <network>\n", utils.GetNodevinExecutable())
	fmt.Printf("%s start <network> --testnet\n", utils.GetNodevinExecutable())
	fmt.Printf("%s start bitcoin --ord\n", utils.GetNodevinExecutable())
	fmt.Printf("%s start litecoin --ord-litecoin\n", utils.GetNodevinExecutable())
}
