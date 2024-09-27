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

package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func mergeConfigs(defaultConfig, overrideConfig NetworkConfig) NetworkConfig {
	if overrideConfig.Image != "" {
		defaultConfig.Image = overrideConfig.Image
	}
	if overrideConfig.Version != "" {
		defaultConfig.Version = overrideConfig.Version
	}
	if overrideConfig.ContainerName != "" {
		defaultConfig.ContainerName = overrideConfig.ContainerName
	}
	if overrideConfig.Restart != "" {
		defaultConfig.Restart = overrideConfig.Restart
	}
	if overrideConfig.Command != "" {
		defaultConfig.Command = overrideConfig.Command
	}
	if len(overrideConfig.Ports) > 0 {
		defaultConfig.Ports = overrideConfig.Ports
	}
	if len(overrideConfig.Volumes) > 0 {
		defaultConfig.Volumes = overrideConfig.Volumes
	}
	if len(overrideConfig.Networks) > 0 {
		defaultConfig.Networks = overrideConfig.Networks
	}
	if overrideConfig.Deploy.Resources.Limits.CPUs != "" {
		defaultConfig.Deploy.Resources.Limits.CPUs = overrideConfig.Deploy.Resources.Limits.CPUs
	}
	if overrideConfig.Deploy.Resources.Limits.Memory != "" {
		defaultConfig.Deploy.Resources.Limits.Memory = overrideConfig.Deploy.Resources.Limits.Memory
	}
	if overrideConfig.Deploy.Resources.Reservations.CPUs != "" {
		defaultConfig.Deploy.Resources.Reservations.CPUs = overrideConfig.Deploy.Resources.Reservations.CPUs
	}
	if overrideConfig.Deploy.Resources.Reservations.Memory != "" {
		defaultConfig.Deploy.Resources.Reservations.Memory = overrideConfig.Deploy.Resources.Reservations.Memory
	}
	if len(overrideConfig.NetworkDefs) > 0 {
		for k, v := range overrideConfig.NetworkDefs {
			defaultConfig.NetworkDefs[k] = v
		}
	}
	if len(overrideConfig.VolumeDefs) > 0 {
		for k, v := range overrideConfig.VolumeDefs {
			defaultConfig.VolumeDefs[k] = v
		}
	}
	return defaultConfig
}

func isDeploySet(deploy Deploy) bool {
	return deploy.Resources.Limits.CPUs != "" ||
		deploy.Resources.Limits.Memory != "" ||
		deploy.Resources.Reservations.CPUs != "" ||
		deploy.Resources.Reservations.Memory != ""
}

func createExtraServices(extraServiceNames []string, extraServiceConfigs []NetworkConfig, extraNetworkDefs map[string]NetworkDetails, extraVolumeDefs map[string]VolumeDetails) (map[string]Service, map[string]NetworkDetails, map[string]VolumeDetails) {
	services := make(map[string]Service)
	networkDefs := make(map[string]NetworkDetails)
	volumeDefs := make(map[string]VolumeDetails)

	for i, serviceName := range extraServiceNames {
		config := extraServiceConfigs[i]

		override := NetworkConfig{
			Image:         viper.GetString(fmt.Sprintf("%s-image", serviceName)),
			Version:       viper.GetString(fmt.Sprintf("%s-version", serviceName)),
			ContainerName: viper.GetString(fmt.Sprintf("%s-container-name", serviceName)),
			Restart:       viper.GetString(fmt.Sprintf("%s-restart", serviceName)),
			Command:       viper.GetString(fmt.Sprintf("%s-command", serviceName)),
			Ports:         viper.GetStringSlice(fmt.Sprintf("%s-ports", serviceName)),
			Volumes:       viper.GetStringSlice(fmt.Sprintf("%s-volumes", serviceName)),
			Networks:      viper.GetStringSlice(fmt.Sprintf("%s-networks", serviceName)),
			Deploy: Deploy{
				Resources: Resources{
					Limits: ResourceDetails{
						CPUs:   viper.GetString(fmt.Sprintf("%s-cpu-limit", serviceName)),
						Memory: viper.GetString(fmt.Sprintf("%s-mem-limit", serviceName)),
					},
					Reservations: ResourceDetails{
						CPUs:   viper.GetString(fmt.Sprintf("%s-cpu-reservation", serviceName)),
						Memory: viper.GetString(fmt.Sprintf("%s-mem-reservation", serviceName)),
					},
				},
			},
			NetworkDefs: extraNetworkDefs,
			VolumeDefs:  extraVolumeDefs,
		}

		finalConfig := mergeConfigs(config, override)

		service := Service{
			Image:         finalConfig.Image + ":" + finalConfig.Version,
			ContainerName: finalConfig.ContainerName,
			Restart:       finalConfig.Restart,
			Command:       finalConfig.Command,
			Ports:         finalConfig.Ports,
			Volumes:       finalConfig.Volumes,
			Networks:      finalConfig.Networks,
		}

		if isDeploySet(finalConfig.Deploy) {
			service.Deploy = &Deploy{Resources: finalConfig.Deploy.Resources}
		}

		services[serviceName] = service
		for k, v := range finalConfig.NetworkDefs {
			networkDefs[k] = v
		}
		for k, v := range finalConfig.VolumeDefs {
			volumeDefs[k] = v
		}
	}

	return services, networkDefs, volumeDefs
}

const GB = 1 << 30 // 1 GB in bytes

func CreateComposeFile(nodeName string, config NetworkConfig, extraServiceNames []string, extraServiceConfigs []NetworkConfig, cwd string) (string, error) {
	nodevinDir, err := GetNodevinDataDir()
	if err != nil {
		return "", fmt.Errorf("failed to create image-specific directory: %w", err)
	}

	// Dynamically generate the sub-directory for this specific image within ~/.nodevin
	err = os.MkdirAll(config.LocalPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create image-specific directory: %w", err)
	}

	// Check if the total size of files in imageDir is greater than 1 GB
	filesNeedCopy := false
	totalSize, err := getDirectorySize(config.LocalPath)
	if err != nil || totalSize < GB {
		filesNeedCopy = true
	}

	// Create the volumes for the services
	dockerNetworks := viper.GetStringSlice("docker-networks")
	volumeDefinitions := viper.GetStringSlice("volume-definitions")

	networkDefs := make(map[string]NetworkDetails)
	volumeDefs := make(map[string]VolumeDetails)

	for _, network := range dockerNetworks {
		if network != "" {
			networkDefs[network] = NetworkDetails{Driver: viper.GetString("network-driver")}
		}
	}

	for _, volume := range volumeDefinitions {
		if volume != "" {
			volumeDefs[volume] = VolumeDetails{Labels: viper.GetStringMapString("volume-labels")}
		}
	}

	override := NetworkConfig{
		Image:         viper.GetString("image"),
		Version:       viper.GetString("version"),
		Restart:       viper.GetString("restart"),
		ContainerName: viper.GetString("container-name"),
		Command:       viper.GetString("command"),
		Ports:         viper.GetStringSlice("ports"),
		Volumes:       viper.GetStringSlice("volumes"),
		Networks:      dockerNetworks,
		Deploy: Deploy{
			Resources: Resources{
				Limits: ResourceDetails{
					CPUs:   viper.GetString("cpu-limit"),
					Memory: viper.GetString("mem-limit"),
				},
				Reservations: ResourceDetails{
					CPUs:   viper.GetString("cpu-reservation"),
					Memory: viper.GetString("mem-reservation"),
				},
			},
		},
		NetworkDefs: networkDefs,
		VolumeDefs:  volumeDefs,
	}

	finalConfig := mergeConfigs(config, override)

	// Main service configuration
	mainService := Service{
		Image:         finalConfig.Image + ":" + finalConfig.Version,
		ContainerName: finalConfig.ContainerName,
		Restart:       finalConfig.Restart,
		Command:       finalConfig.Command,
		Ports:         finalConfig.Ports,
		Volumes:       finalConfig.Volumes,
		Networks:      finalConfig.Networks,
	}

	// Initialize services map
	services := make(map[string]Service)

	// Add init container service only if files need to be copied
	if filesNeedCopy {
		initContainerName := fmt.Sprintf("init-config-%s", nodeName)
		initVolumeName := fmt.Sprintf("%s-init-volume", nodeName)

		initService := Service{
			Image:         finalConfig.Image + ":" + finalConfig.Version,
			ContainerName: initContainerName,
			Restart:       "no",
			Command: fmt.Sprintf(`/bin/sh -c "
if [ -z \"$(ls -A /nodevin-volume)\" ]; then
  mkdir -p /nodevin-volume/ &&
  cp -r * /nodevin-volume/ &&
  touch /nodevin-volume/.copy-done
else
  echo 'Volume not empty, skipping file copy';
  touch /nodevin-volume/.copy-done;
fi"`),
			Volumes: []string{
				fmt.Sprintf("%s:/init-volume", initVolumeName),
				fmt.Sprintf("%s:/nodevin-volume", config.LocalPath),
			},
			Entrypoint: "",
		}

		// Add init container service to the services map
		services[initContainerName] = initService

		// Add dependency for mainService to wait for initService to complete
		mainService.DependsOn = map[string]ServiceDependsOnCondition{
			initContainerName: {
				Condition: "service_completed_successfully",
			},
		}
	}

	services[nodeName] = mainService

	// Include any extra services
	extraNetworkDefs := finalConfig.NetworkDefs
	extraVolumeDefs := finalConfig.VolumeDefs

	if len(extraServiceNames) > 0 && len(extraServiceConfigs) > 0 {
		extraServices, extraNetworks, extraVolumes := createExtraServices(extraServiceNames, extraServiceConfigs, extraNetworkDefs, extraVolumeDefs)
		for k, v := range extraServices {
			services[k] = v
		}

		for k, v := range extraNetworks {
			extraNetworkDefs[k] = v
		}
		for k, v := range extraVolumes {
			extraVolumeDefs[k] = v
		}
	}

	// Build the compose file structure
	composeFile := ComposeFile{
		Version:  "3.9",
		Services: services,
		Networks: extraNetworkDefs,
		Volumes: map[string]VolumeDetails{
			fmt.Sprintf("%s-init-volume", nodeName): {
				Labels: map[string]string{
					"nodevin.init.volume":         "true",
					"nodevin.blockchain.software": fmt.Sprintf("%s-init-volume", nodeName),
				},
			},
		},
	}

	// Generate and save the Compose file
	composeFileName := fmt.Sprintf("docker-compose_%s.yml", nodeName)
	composeFilePath := filepath.Join(nodevinDir, composeFileName)
	composeData, err := yaml.Marshal(&composeFile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal docker-compose.yml: %w", err)
	}

	if err = os.WriteFile(composeFilePath, composeData, 0644); err != nil {
		return "", fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	return composeFilePath, nil
}

// Helper function to calculate the total size of files in a directory
func getDirectorySize(dir string) (int64, error) {
	var totalSize int64 = 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Accumulate file size
			totalSize += info.Size()
		}
		return nil
	})
	return totalSize, err
}
