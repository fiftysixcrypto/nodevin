package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version of nodevin",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Nodevin CLI v0.1.0")
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}

