package deploy

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/melbahja/goph"
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

func uploadTextFile(c *goph.Client, filename, content string) error {
	sftp, err := c.NewSftp()
	if err != nil {
		return err
	}
	f, err := sftp.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(content))
	return err
}
