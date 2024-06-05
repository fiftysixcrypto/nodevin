package utils

import (
	"sort"
	"strings"
)

var networkContainerMap = map[string]string{
	"bitcoin":         "bitcoin-core",
	"ethereum":        "geth",
	"ethereumclassic": "core-geth",
	"litecoin":        "litecoin-core",
	"dogecoin":        "dogecoin-core",
}

func NetworkContainerMap() map[string]string {
	return networkContainerMap
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

func CheckIfTestnetOrTestnetNetwork() bool {
	return true
}
