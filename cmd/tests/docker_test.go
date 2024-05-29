// cmd/tests/docker_test.go
package cmd

import (
    "context"
    "testing"

    "github.com/docker/docker/api/types/filters"
    "github.com/docker/docker/client"
    "github.com/stretchr/testify/assert"
)

// Mock docker client for testing
var mockDockerClient *client.Client

func initMockDockerClient() error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }
    mockDockerClient = cli
    return nil
}

func createTestVolume(volumeName string) error {
    ctx := context.Background()
    _, err := mockDockerClient.VolumeCreate(ctx, volume.CreateOptions{
        Name: volumeName,
    })
    return err
}

func removeTestVolume(volumeName string) error {
    ctx := context.Background()
    return mockDockerClient.VolumeRemove(ctx, volumeName, true)
}

func TestCreateVolume(t *testing.T) {
    err := initMockDockerClient()
    assert.NoError(t, err)

    err = createTestVolume("testvolume")
    assert.NoError(t, err)

    volumes, err := mockDockerClient.VolumeList(context.Background(), filters.Args{})
    assert.NoError(t, err)

    found := false
    for _, v := range volumes.Volumes {
        if v.Name == "testvolume" {
            found = true
            break
        }
    }
    assert.True(t, found)

    err = removeTestVolume("testvolume")
    assert.NoError(t, err)
}

