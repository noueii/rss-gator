package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, configFileName)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaulCfg := &Config{
			CurrentUsername: "",
			DbURL:           "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
		}

		return defaulCfg, defaulCfg.write()

	}

	content, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(content, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) SetUser(username string) error {
	initial := *c
	c.CurrentUsername = username
	err := c.write()
	if err != nil {
		c = &initial
		return err
	}
	return nil
}

func (c *Config) write() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, configFileName)

	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
