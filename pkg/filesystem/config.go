package filesystem

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Preferences struct {
	AutoUpdate  bool `json:"auto_update"`
	AutoUpgrade bool `json:"auto_upgrade"`
}

type Configuration struct {
	SchemaVersion string `json:"schema_version"`
	Developer     string `json:"developer_name"`

	GitHubToken string      `json:"github_token,omitempty"`
	GitlabToken string      `json:"gitlab_key,omitempty"`
	Preferences Preferences `json:"preferences"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".chiroptera"), nil
}

func SaveConfiguration(configuration *Configuration) error {
	directory, err := getConfigPath()
	configPath := filepath.Join(directory, "config.json")
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(configuration, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

type GetKeyMode int

const (
	GetPrivateKey GetKeyMode = iota
	GetPublicKey
	GetBoth
)

func GetKeys(mode GetKeyMode) ([]byte, []byte, error) {
	directory, err := getConfigPath()
	if err != nil {
		return nil, nil, err
	}

	var priv, pub []byte

	if mode == GetPrivateKey || mode == GetBoth {
		privKeyFile := filepath.Join(directory, "keys", "key")
		priv, err = os.ReadFile(privKeyFile)
		if err != nil {
			return nil, nil, err
		}
	}

	if mode == GetPublicKey || mode == GetBoth {
		pubKeyFile := filepath.Join(directory, "keys", "key.pub")
		pub, err = os.ReadFile(pubKeyFile)
		if err != nil {
			return nil, nil, err
		}
	}

	return priv, pub, nil
}

func LoadConfiguration() (*Configuration, error) {
	directory, err := getConfigPath()
	configPath := filepath.Join(directory, "config.json")
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var configuration Configuration
	if err := json.Unmarshal(configData, &configuration); err != nil {
		return nil, err
	}

	return &configuration, nil
}

func CheckConfiguration() bool {
	_, err := LoadConfiguration()
	if err != nil {
		return false
	}

	_, _, err = GetKeys(GetBoth)
	if err != nil {
		return false
	}

	return true
}

func PurgeFolder() error {
	directory, err := getConfigPath()
	if err != nil {
		return err
	}

	if err = os.RemoveAll(directory); err != nil {
		return err
	}

	return nil
}