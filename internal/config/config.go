package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey     string `json:"api_key"`
	APIURL     string `json:"api_url"`
	LogEnabled bool   `json:"log_enabled,omitempty"`
}

const defaultAPIURL = "https://battlefaeries.vercel.app"

func configDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".battlefaeries")
}

func configPath() string {
	return filepath.Join(configDir(), "config.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{APIURL: defaultAPIURL}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.APIURL == "" {
		cfg.APIURL = defaultAPIURL
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	if err := os.MkdirAll(configDir(), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath(), data, 0600)
}

func LogDir() string {
	return filepath.Join(configDir(), "logs")
}

func LogPath() string {
	return filepath.Join(LogDir(), "activity.jsonl")
}
