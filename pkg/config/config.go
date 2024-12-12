package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
)

type Config struct {
	DownloadsHomeDir string `json:"downloads_home_dir"`
	IaKeyFile        string `json:"ia_key_file"`
	CookiesFile      string `json:"cookies_file"`
}

func NewConfig() *Config {
	iaKeyFile := ".aixifan_ia_keys.txt"
	cookiesFile := ".cookies.txt"
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Info("Failed to get user home dir, using current dir instead")
		home = "."
	}
	iaKeyFile = path.Join(home, iaKeyFile)
	cookiesFile = path.Join(home, cookiesFile)

	return &Config{
		DownloadsHomeDir: "aixifan_downloads",
		IaKeyFile:        iaKeyFile,
		CookiesFile:      cookiesFile,
	}
}

func (config *Config) Save() error {
	json_data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile("aixifan_config.json", json_data, 0644)
	return err
}

func LoadConfig() (*Config, error) {
	json_data, err := os.ReadFile("aixifan_config.json")
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(json_data, &config)
	return &config, err
}

func LoadOrNewConfig() (*Config, error) {
	config, err := LoadConfig()
	if err != nil {
		slog.Info("Failed to load config, creating new one")
		config = NewConfig()
		err = config.Save()
	}
	return config, err
}

func (c *Config) MakeDownloadsHomeDir() error {
	err := os.MkdirAll(c.DownloadsHomeDir, 0755)
	return err
}
