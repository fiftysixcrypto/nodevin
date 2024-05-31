package env

import (
	"fmt"
	"os"
)

func WriteEnvVariable(file *os.File, key, value string) error {
	if len(value) > 0 {
		_, err := fmt.Fprintf(file, "%s=%s\n", key, value)
		if err != nil {
			return fmt.Errorf("failed to write %s to .env file: %w", key, err)
		}
	}
	return nil
}
