package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/curveballdaniel/nodevin/internal/logger"
)

const (
	updateURL      = "https://github.com/fiftysixcrypto/nodevin/releases/latest/download/nodevin"
	currentVersion = "1.0.0" // Replace with the current version of your application
)

func CheckForUpdates() (bool, error) {
	resp, err := http.Get("https://api.github.com/repos/fiftysixcrypto/nodevin/releases/latest")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to fetch update info: %s", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, err
	}

	if release.TagName == currentVersion {
		return false, nil
	}

	logger.LogInfo(fmt.Sprintf("New version available: %s", release.TagName))
	if err := downloadUpdate(); err != nil {
		return false, fmt.Errorf("failed to download update: %w", err)
	}

	return true, nil
}

func downloadUpdate() error {
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

func ApplyUpdate() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "move /Y nodevin_new nodevin.exe")
	} else {
		cmd = exec.Command("sh", "-c", "mv nodevin_new nodevin && chmod +x nodevin")
	}
	return cmd.Run()
}