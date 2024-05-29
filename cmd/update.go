package cmd

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "runtime"

    "github.com/spf13/cobra"
)

const (
    updateURL = "https://github.com/yourusername/nodevin/releases/latest/download/nodevin"
)

func checkForUpdates() error {
    resp, err := http.Get(updateURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to fetch update: %s", resp.Status)
    }

    out, err := os.Create("nodevin_new")
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return err
    }

    return nil
}

func applyUpdate() error {
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/C", "move /Y nodevin_new nodevin.exe")
    } else {
        cmd = exec.Command("sh", "-c", "mv nodevin_new nodevin && chmod +x nodevin")
    }
    return cmd.Run()
}

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Check for updates and apply them if available",
    Run: func(cmd *cobra.Command, args []string) {
        logInfo("Checking for updates...")
        if err := checkForUpdates(); err != nil {
            logError("Failed to check for updates: " + err.Error())
            return
        }
        logInfo("Update downloaded. Applying update...")
        if err := applyUpdate(); err != nil {
            logError("Failed to apply update: " + err.Error())
            return
        }
        logInfo("Update applied successfully. Please restart the application.")
    },
}

func init() {
    rootCmd.AddCommand(updateCmd)
}

