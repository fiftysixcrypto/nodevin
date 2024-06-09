package utils

import (
	"sort"
	"strings"

	"github.com/spf13/viper"
)

type NetworkInfo struct {
	ContainerName string
	RPCPort       int
	StartMessage  string
}

var networkInfoMap = map[string]NetworkInfo{
	"bitcoin": {
		ContainerName: "bitcoin-core",
		RPCPort:       8332,
		StartMessage:  "\"A system for electronic transactions without relying on trust.\" -- Satoshi Nakamoto",
	},
	"bitcoin-testnet": {
		ContainerName: "bitcoin-core-testnet",
		RPCPort:       18332,
		StartMessage:  "\"Testing is the lifeblood of innovation and security.\"",
	},
	"ethereum": {
		ContainerName: "geth",
		RPCPort:       8545,
		StartMessage:  "\"Ethereum is the foundation for our digital future.\" -- Vitalik Buterin",
	},
	"ethereum-classic": {
		ContainerName: "core-geth",
		RPCPort:       8545,
		StartMessage:  "\"Code is law.\" -- Ethereum Classic",
	},
	"litecoin": {
		ContainerName: "litecoin-core",
		RPCPort:       9332,
		StartMessage:  "\"Litecoin is the silver to Bitcoin's gold.\" -- Charlie Lee",
	},
	"dogecoin": {
		ContainerName: "dogecoin-core",
		RPCPort:       22555,
		StartMessage:  "\"To the moon!\" -- Dogecoin Community",
	},
}

func NetworkContainerMap() map[string]string {
	networkContainerMap := make(map[string]string)
	for key, value := range networkInfoMap {
		networkContainerMap[key] = value.ContainerName
	}
	return networkContainerMap
}

func NetworkDefaultRPCPorts() map[string]int {
	networkDefaultRPCPorts := make(map[string]int)
	for key, value := range networkInfoMap {
		networkDefaultRPCPorts[key] = value.RPCPort
	}
	return networkDefaultRPCPorts
}

func GetStartMessage(network string) (string, bool) {
	networkInfo, exists := networkInfoMap[network]
	return networkInfo.StartMessage, exists
}

func GetFiftysixLocalMappedContainerName(network string) (string, bool) {
	networkInfo, exists := networkInfoMap[network]
	return networkInfo.ContainerName, exists
}

func GetFiftysixDockerhubContainerName(network string) (string, bool) {
	networkInfo, exists := networkInfoMap[network]
	return "fiftysix/" + networkInfo.ContainerName, exists
}

func GetAllSupportedNetworks() string {
	var keys []string
	for key := range networkInfoMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return strings.Join(keys, ", ")
}

func CheckIfTestnetOrTestnetNetworkFlag() bool {
	networkFlag := viper.GetString("network")
	testnetFlag := viper.GetBool("testnet")

	return testnetFlag || networkFlag == "testnet"
}
