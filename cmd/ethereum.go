package cmd

import (
    "fmt"
)

type Ethereum struct{}

func (e Ethereum) StartNode(config Config) error {
    fmt.Println("Starting Ethereum node with config:", config)
    // Add logic to start Ethereum node
    return nil
}

func (e Ethereum) StopNode() error {
    fmt.Println("Stopping Ethereum node")
    // Add logic to stop Ethereum node
    return nil
}

func init() {
    RegisterBlockchain("ethereum", Ethereum{})
}

