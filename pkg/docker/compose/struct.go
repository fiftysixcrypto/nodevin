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
	Version       string
	ContainerName string
	Command       string
	Ports         []string
	Volumes       []string
	Networks      []string
	Deploy        Deploy
	NetworkDefs   map[string]NetworkDetails
	VolumeDefs    map[string]VolumeDetails
}
