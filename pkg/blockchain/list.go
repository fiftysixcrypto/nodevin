package blockchain

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fiftysixcrypto/nodevin/internal/utils"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all supported blockchain networks",
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

	fmt.Printf("Supported networks: %s\n", strings.Join(networkNames, ", "))
}
