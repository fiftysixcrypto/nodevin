// cmd/tests/integration_test.go
package cmd

import (
    "context"
    "os/exec"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestDaemonAndUpdate(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    cmd := exec.CommandContext(ctx, "../nodevin", "daemon", "--network", "bitcoin", "--storage_path", "/tmp/bitcoin", "--port", "40404", "--resource_limit", "1GB")
    err := cmd.Start()
    assert.NoError(t, err)

    // Allow some time for the daemon to start
    time.Sleep(5 * time.Second)

    // Check for updates
    err = checkForUpdates()
    assert.NoError(t, err)

    // Stop the daemon
    cancel()
    err = cmd.Wait()
    assert.NoError(t, err)
}

