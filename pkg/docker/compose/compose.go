package compose

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func mergeConfigs(defaultConfig, overrideConfig NetworkConfig) NetworkConfig {
	if overrideConfig.Image != "" {
		defaultConfig.Image = overrideConfig.Image
	}
	if overrideConfig.ContainerName != "" {
		defaultConfig.ContainerName = overrideConfig.ContainerName
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

func CreateComposeFile(nodeName string, cwd string, config NetworkConfig, overrides NetworkConfig) (string, error) {
	// Merge default and override configurations
	finalConfig := mergeConfigs(config, overrides)

	// Define the Docker Compose structure
	service := Service{
		Image:         finalConfig.Image,
		ContainerName: finalConfig.ContainerName,
		Command:       finalConfig.Command,
		Ports:         finalConfig.Ports,
		Volumes:       finalConfig.Volumes,
		Networks:      finalConfig.Networks,
	}

	if isDeploySet(finalConfig.Deploy) {
		// Only set non-empty fields in the deploy section
		deploy := &Deploy{}
		if finalConfig.Deploy.Resources.Limits.CPUs != "" || finalConfig.Deploy.Resources.Limits.Memory != "" {
			deploy.Resources.Limits = finalConfig.Deploy.Resources.Limits
		}
		if finalConfig.Deploy.Resources.Reservations.CPUs != "" || finalConfig.Deploy.Resources.Reservations.Memory != "" {
			deploy.Resources.Reservations = finalConfig.Deploy.Resources.Reservations
		}
		service.Deploy = deploy
	}

	composeFile := ComposeFile{
		Version: "3.9",
		Services: map[string]Service{
			"main-service": service,
		},
		Networks: finalConfig.NetworkDefs,
		Volumes:  finalConfig.VolumeDefs,
	}

	// Write the Docker Compose file to the current working directory
	composeFileName := fmt.Sprintf("docker-compose_%s.yml", nodeName)
	composeFilePath := filepath.Join(cwd, composeFileName)
	composeData, err := yaml.Marshal(&composeFile)
	if err != nil {
		return "", fmt.Errorf("failed to marshal docker-compose.yml file: %w", err)
	}

	err = os.WriteFile(composeFilePath, composeData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write docker-compose.yml file: %w", err)
	}

	return composeFilePath, nil
}
