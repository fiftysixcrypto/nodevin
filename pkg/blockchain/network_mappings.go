package blockchain

var networkContainerMap = map[string]string{
	"bitcoin":  "bitcoin-core",
	"ethereum": "geth",
	// Add more mappings as needed
}

func getFiftysixLocalMappedContainerName(network string) (string, bool) {
	containerName, exists := networkContainerMap[network]
	return containerName, exists
}

func getFiftysixDockerhubContainerName(network string) (string, bool) {
	containerName, exists := networkContainerMap[network]
	return "fiftysix/" + containerName, exists
}

/*
var networkContainerMap = map[string]string{
    "bitcoin":         "bitcoin-core",
    "ethereum":        "geth",
    "ethereumclassic": "etc-geth",
    "litecoin":        "litecoin-core",
    "dogecoin":        "dogecoin-core",
}*/
