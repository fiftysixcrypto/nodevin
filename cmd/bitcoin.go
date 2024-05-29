package cmd

import (
    "fmt"
)

type Bitcoin struct{}

func (b Bitcoin) StartNode(config Config) error {
    fmt.Println("Starting Bitcoin node with config:", config)
    // Add logic to start Bitcoin node
    return nil
}

func (b Bitcoin) StopNode() error {
    fmt.Println("Stopping Bitcoin node")
    // Add logic to stop Bitcoin node
    return nil
}

func init() {
    RegisterBlockchain("bitcoin", Bitcoin{})
}

