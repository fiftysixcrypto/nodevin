/*
// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2024 The Nodevin Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
*/

package initialize

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"

	"github.com/fiftysixcrypto/nodevin/internal/logger"
	"github.com/fiftysixcrypto/nodevin/internal/version"
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize nodevin and check system capabilities",
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func runInit() {
	fmt.Println("")
	fmt.Printf("Welcome to nodevin v%s!\n", version.Version)
	fmt.Println("Nodevin allows you to run any blockchain node in the world.")
	fmt.Println("")
	fmt.Println("--")
	fmt.Println("")

	// inspection
	fmt.Println("Nodevin will now inspect your system for docker and docker compose versions...\n")
	if err := performInspection(); err != nil {
		fmt.Println("")
		logger.LogError("System inspection failed: " + err.Error())
		if err := installDockerAndCompose(); err != nil {
			logger.LogError("Failed to install Docker and Docker Compose: " + err.Error())
			return
		}
	}

	fmt.Println("")
	fmt.Println("--")
	fmt.Println("")

	// outro
	fmt.Println("Thank you for using nodevin!")
	fmt.Println("It's time to start your own Bitcoin node. Run `nodevin start bitcoin` to get started.")
}

func performInspection() error {
	fmt.Print("Checking Docker version... ")
	if err := checkDockerVersion(); err != nil {
		return err
	}

	fmt.Print("Checking Docker Compose version... ")
	if err := checkDockerComposeVersion(); err != nil {
		return err
	}

	return nil
}

func checkDockerVersion() error {
	cmd := exec.Command("docker", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check Docker version: %w", err)
	}

	versionInfo := string(output)
	fmt.Print(versionInfo)

	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(versionInfo)
	if len(matches) < 4 {
		return fmt.Errorf("unable to parse Docker version")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	if major < 20 {
		return fmt.Errorf("Docker version 20+ is required. Please update Docker.")
	}

	fmt.Printf("Docker version %d.%d.%d meets the requirement.\n", major, minor, patch)
	return nil
}

func checkDockerComposeVersion() error {
	cmd := exec.Command("docker-compose", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to check Docker Compose version: %w", err)
	}

	versionInfo := string(output)
	fmt.Print(versionInfo)

	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(versionInfo)
	if len(matches) < 4 {
		return fmt.Errorf("unable to parse Docker Compose version")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	if major < 1 {
		return fmt.Errorf("Docker Compose version 1+ is required. Please update Docker Compose.")
	}

	fmt.Printf("Docker Compose version %d.%d.%d meets the requirement.\n", major, minor, patch)
	return nil
}

func installDockerAndCompose() error {
	fmt.Println("Installing Docker and Docker Compose...")

	if err := installDocker(); err != nil {
		return fmt.Errorf("failed to install Docker: %w", err)
	}

	if err := installDockerCompose(); err != nil {
		return fmt.Errorf("failed to install Docker Compose: %w", err)
	}

	fmt.Println("Docker and Docker Compose installed successfully.")
	return nil
}

func installDocker() error {
	var installCmd *exec.Cmd

	fmt.Println("Installing Docker...")

	switch runtime.GOOS {
	case "linux":
		installCmd = exec.Command("sh", "-c", `
			sudo apt update &&
			sudo apt install -y apt-transport-https ca-certificates curl software-properties-common &&
			curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - &&
			sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" &&
			sudo apt update &&
			sudo apt install -y docker-ce`)
	case "darwin":
		installCmd = exec.Command("sh", "-c", `
			brew install --cask docker &&
			open /Applications/Docker.app &&
			while ! docker system info > /dev/null 2>&1; do sleep 1; done`)
	case "windows":
		installCmd = exec.Command("powershell", "-Command", `
			$ErrorActionPreference = 'Stop'; 
			Invoke-WebRequest -UseBasicParsing -OutFile docker-desktop-installer.exe https://desktop.docker.com/win/stable/Docker%20Desktop%20Installer.exe;
			Start-Process -FilePath docker-desktop-installer.exe -Wait -ArgumentList @("--quiet");
			Remove-Item -Force docker-desktop-installer.exe;
			Start-Process "Docker Desktop" -Wait;
			while ((docker version) -eq $null) { Start-Sleep -Seconds 5 }`)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err := installCmd.Run(); err != nil {
		return err
	}

	return nil
}

func installDockerCompose() error {
	var installCmd *exec.Cmd

	fmt.Println("Installing Docker Compose...")

	switch runtime.GOOS {
	case "linux":
		installCmd = exec.Command("sh", "-c", `
			sudo curl -L "https://github.com/docker/compose/releases/download/v2.27.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose &&
			sudo chmod +x /usr/local/bin/docker-compose &&
			docker-compose --version`)
	case "darwin":
		installCmd = exec.Command("sh", "-c", "brew install docker-compose")
	case "windows":
		installCmd = exec.Command("powershell", "-Command", `
			$ErrorActionPreference = 'Stop'; 
			Invoke-WebRequest "https://github.com/docker/compose/releases/download/v2.27.1/docker-compose-Windows-x86_64.exe" -OutFile "$env:ProgramFiles\Docker\Docker\resources\bin\docker-compose.exe";
			& "$env:ProgramFiles\Docker\Docker\resources\bin\docker-compose.exe" --version`)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err := installCmd.Run(); err != nil {
		return err
	}

	return nil
}