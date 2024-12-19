package deploy

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

func getDefaultSSHPrivateKeyPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	homeDir := currentUser.HomeDir
	if homeDir == "" {
		return "", fmt.Errorf("can't find home folder")
	}

	defaultKeys := []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
	}

	for _, keyPath := range defaultKeys {
		if _, err := os.Stat(keyPath); err == nil {
			return keyPath, nil
		}
	}

	return "", fmt.Errorf("can't find default key")
}
