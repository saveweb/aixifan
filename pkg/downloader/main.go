package downloader

import (
	"flag"
	"log/slog"

	"github.com/saveweb/aixifan/pkg/config"
	"github.com/saveweb/aixifan/pkg/utils"
)

func Main(downCmd *flag.FlagSet, downNoVersionCheck *bool, downDougaId *string, downSkipIACheck *bool) int {
	dougaId := *downDougaId

	// check new version
	newVersionChan := make(chan bool, 1)
	defer func(ch chan bool) {
		select {
		case newVersion := <-ch:
			if newVersion {
				slog.Warn("")
				slog.Warn("New version available, please update :) <https://github.com/saveweb/aixifan/releases/latest>")
			}
		default:
		}
	}(newVersionChan)
	defer close(newVersionChan)
	if !*downNoVersionCheck {
		go func(ch chan bool) {
			slog.Info("Checking new version...")
			newVersion, err := utils.NewVersionAvailable()
			if err != nil {
				slog.Warn("Failed to check new version", "err", err)
			}
			if newVersion {
				slog.Warn("New version available, please update :) <https://github.com/saveweb/aixifan/releases/latest>")
				ch <- true
			} else {
				slog.Info("You are using the latest version :)")
				ch <- false
			}
		}(newVersionChan)
	}

	// Check IA Item (dedup)
	if !*downSkipIACheck {
		slog.Info("Checking IA item...")
		identifier := utils.ToIdentifier(dougaId)
		exists, err := utils.CheckIAItemExist(identifier)
		if err != nil {
			slog.Error("Failed to check IA item", "err", err)
			return 1
		}
		if exists {
			slog.Info("Item already exists on IA", "identifier", identifier)
			return 88
		}
		slog.Info("Item does not exist on IA", "identifier", identifier)
	}

	// download
	if dougaId == "" {
		downCmd.Usage()
		return 2
	}

	config, err := config.LoadOrNewConfig()
	if err != nil {
		panic(err)
	}
	if err := config.MakeDownloadsHomeDir(); err != nil {
		panic(err)
	}

	if err := Download(config.DownloadsHomeDir, dougaId); err != nil {
		slog.Error("Failed to download", "err", err)
		return 1
	}

	return 0
}
