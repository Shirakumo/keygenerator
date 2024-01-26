package config

import (
	"keygenerator/keygen"
	"path/filepath"
	"encoding/json"
	"os"
)

type Config struct {
	KeyURL string `json:"key-url"`
	LastChecked int64 `json:"last-checked"`
	LocalFile *keygen.File `json:"local-file"`
	LocalPath string `json:"local-path"`
}

func Default(keyURL string) Config {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }

	var config Config
	config.KeyURL = keyURL
	config.LastChecked = 0
	config.LocalPath = filepath.Dir(ex)
	return config
}

func Read(path string) (Config, error) {
	var config = Default("")

	b, err := os.ReadFile(path)
    if err != nil {
		return config, err
    }

	err = json.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}
	
	return config, nil
}

func Write(config Config, path string) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0660)
}

func defaultPath() string {
	ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
	return filepath.Join(filepath.Dir(ex), ".key")
}

func ReadDefault() (Config, error) {
	return Read(defaultPath())
}

func WriteDefault(config Config) error {
	return Write(config, defaultPath())
}
