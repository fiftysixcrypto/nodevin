package cmd

type Blockchain interface {
    StartNode(config Config) error
    StopNode() error
}

var blockchainRegistry = make(map[string]Blockchain)

func RegisterBlockchain(name string, blockchain Blockchain) {
    blockchainRegistry[name] = blockchain
}

func GetBlockchain(name string) (Blockchain, bool) {
    blockchain, exists := blockchainRegistry[name]
    return blockchain, exists
}

