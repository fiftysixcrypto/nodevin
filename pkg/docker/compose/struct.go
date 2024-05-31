package compose

type ResourceDetails struct {
	CPUs   string `yaml:"cpus,omitempty"`
	Memory string `yaml:"memory,omitempty"`
}

type Resources struct {
	Limits       ResourceDetails `yaml:"limits,omitempty"`
	Reservations ResourceDetails `yaml:"reservations,omitempty"`
}

type Deploy struct {
	Resources Resources `yaml:"resources,omitempty"`
}

type Service struct {
	Image         string   `yaml:"image"`
	ContainerName string   `yaml:"container_name"`
	Command       string   `yaml:"command"`
	Ports         []string `yaml:"ports"`
	Volumes       []string `yaml:"volumes"`
	Networks      []string `yaml:"networks"`
	Deploy        *Deploy  `yaml:"deploy,omitempty"`
}

type NetworkDetails struct {
	Driver string `yaml:"driver"`
}

type VolumeDetails struct {
	Labels map[string]string `yaml:"labels"`
}

type ComposeFile struct {
	Version  string                    `yaml:"version"`
	Services map[string]Service        `yaml:"services"`
	Networks map[string]NetworkDetails `yaml:"networks"`
	Volumes  map[string]VolumeDetails  `yaml:"volumes"`
}

type NetworkConfig struct {
	Image         string
	ContainerName string
	Command       string
	Ports         []string
	Volumes       []string
	Networks      []string
	Deploy        Deploy
	NetworkDefs   map[string]NetworkDetails
	VolumeDefs    map[string]VolumeDetails
}
