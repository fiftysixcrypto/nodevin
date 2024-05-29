package cmd

var networkContainerMap = map[string]string{
	"bitcoin":  "bitcoin-core",
	"ethereum": "geth",
	// Add more mappings as needed
}

func getContainerName(network string) (string, bool) {
	containerName, exists := networkContainerMap[network]
	return containerName, exists
}
