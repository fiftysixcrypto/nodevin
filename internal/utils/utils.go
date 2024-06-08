package utils

import (
	"sort"
	"strings"

	"github.com/spf13/viper"
)

var networkContainerMap = map[string]string{
	"bitcoin":          "bitcoin-core",
	"bitcoin-testnet":  "bitcoin-core-testnet",
	"ethereum":         "geth",
	"ethereum-classic": "core-geth",
	"litecoin":         "litecoin-core",
	"dogecoin":         "dogecoin-core",
}

var networkDefaultRPCPorts = map[string]int{
	"bitcoin":         8332,
	"bitcoin-testnet": 18332,
	"ethereum":        8545,
}

func NetworkContainerMap() map[string]string {
	return networkContainerMap
}

func NetworkDefaultRPCPorts() map[string]int {
	return networkDefaultRPCPorts
}

func GetFiftysixLocalMappedContainerName(network string) (string, bool) {
	containerName, exists := networkContainerMap[network]
	return containerName, exists
}

func GetFiftysixDockerhubContainerName(network string) (string, bool) {
	containerName, exists := networkContainerMap[network]
	return "fiftysix/" + containerName, exists
}

func GetAllSupportedNetworks() string {
	var keys []string
	for key := range networkContainerMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return strings.Join(keys, ", ")
}

func CheckIfTestnetOrTestnetNetworkFlag() bool {
	networkFlag := viper.GetString("network")
	testnetFlag := viper.GetBool("testnet")

	if testnetFlag || networkFlag == "testnet" {
		return true
	}

	return false
}
