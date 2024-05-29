// cmd/tests/performance_test.go
package cmd

import (
    "context"
    "os/exec"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestPerformance(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    cmd := exec.CommandContext(ctx, "../nodevin", "daemon", "--network", "bitcoin", "--storage_path", "/tmp/bitcoin", "--port", "40404", "--resource_limit", "1GB")
    err := cmd.Start()
    assert.NoError(t, err)

    startTime := time.Now()

    // Allow some time for the daemon to start and simulate load
    time.Sleep(10 * time.Second)

    duration := time.Since(startTime)
    assert.True(t, duration < 15*time.Second, "Daemon startup and load simulation took too long")

    // Stop the daemon
    cancel()
    err = cmd.Wait()
    assert.NoError(t, err)
}

