package config

import (
	"os"
	"testing"
)

func Test_Config(t *testing.T) {
	config := NewConfig()
	config.Save()
	defer os.Remove("aixifan_config.json")

	if _, err := os.Stat("aixifan_config.json"); err != nil {
		t.Fatal(err)
		return
	}

	config_loaded, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
		return
	}

	if config_loaded.DownloadsHomeDir != config.DownloadsHomeDir {
		t.Fatal("DownloadsHomeDir not match")
	}
}
