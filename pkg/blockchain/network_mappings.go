package blockchain

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

func getFiftysixLocalMappedContainerName(network string) (string, bool) {
	containerName, exists := networkContainerMap[network]
	return containerName, exists
}

func getFiftysixDockerhubContainerName(network string) (string, bool) {
	containerName, exists := networkContainerMap[network]
	return "fiftysix/" + containerName, exists
}

func getAllSupportedNetworks() string {
	var keys []string
	for key := range networkContainerMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return strings.Join(keys, ", ")
}
