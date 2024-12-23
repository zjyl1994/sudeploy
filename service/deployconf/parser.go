package deployconf

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/zjyl1994/sudeploy/infra/typedef"
)

func Load(filename string) (*typedef.DeployConf, error) {
	bConf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var result typedef.DeployConf
	err = json.Unmarshal(bConf, &result)
	if err != nil {
		return nil, err
	}
	result.Name = strcase.ToSnake(result.Name)
	if result.User == "" {
		u, err := user.Current()
		if err != nil {
			return nil, err
		}
		result.User = u.Username
	}
	if result.Key == "" {
		defaultKey, err := getDefaultSSHPrivateKeyPath()
		if err != nil {
			return nil, err
		}
		result.Key = defaultKey
	}
	return &result, nil
}

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
