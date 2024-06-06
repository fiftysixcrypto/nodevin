package blockchain

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/curveballdaniel/nodevin/internal/utils"

	"github.com/spf13/cobra"
)

var (
	detach     bool
	user       string
	workdir    string
	envVars    []string
	envFile    string
	privileged bool
)

var shellCmd = &cobra.Command{
	Use:   "shell [network]",
	Short: "Run a shell in the specified blockchain node container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		network := args[0]
		containerName, exists := utils.GetFiftysixLocalMappedContainerName(network)
		if !exists {
			fmt.Printf("Unsupported blockchain network: %s\n", network)
			return
		}
		runShell(containerName)
	},
}

func runShell(containerName string) {
	args := []string{"exec", "-it"}
	if detach {
		args = append(args, "-d")
	}
	if user != "" {
		args = append(args, "-u", user)
	}
	if workdir != "" {
		args = append(args, "-w", workdir)
	}
	for _, env := range envVars {
		args = append(args, "-e", env)
	}
	if envFile != "" {
		args = append(args, "--env-file", envFile)
	}
	if privileged {
		args = append(args, "--privileged")
	}
	args = append(args, containerName, "/bin/bash")

	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to run shell in container %s: %v\n", containerName, err)
	}
}

func init() {
	shellCmd.Flags().BoolVarP(&detach, "detach", "d", false, "Run command in the background")
	shellCmd.Flags().StringVarP(&user, "user", "u", "", "Username or UID to run the command as")
	shellCmd.Flags().StringVarP(&workdir, "workdir", "w", "", "Working directory inside the container")
	shellCmd.Flags().StringArrayVarP(&envVars, "env", "e", nil, "Set environment variables")
	shellCmd.Flags().StringVar(&envFile, "env-file", "", "Read environment variables from a file")
	shellCmd.Flags().BoolVar(&privileged, "privileged", false, "Give extended privileges to the command")
}