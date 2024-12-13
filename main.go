package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/saveweb/aixifan/pkg/config"
	"github.com/saveweb/aixifan/pkg/downloader"
	"github.com/saveweb/aixifan/pkg/utils"
)

func rootHelp() {
	fmt.Println("aixifan - " + utils.GetVersion().Version + " (" + utils.GetVersion().GoVersion + ")")
	fmt.Println("")
	fmt.Println(
		`Usage: aixifan COMMAND [ARGS]...

Commands:
  down  Download
  up    Upload to IA
  version  Show version`)
}

func main() {
	flag.Usage = rootHelp
	flag.Parse()

	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	downNoVersionCheck := downCmd.Bool("no-version-check", false, "Do not check aixifan's version")
	downDougaId := downCmd.String("i", "", "Douga ID (int-string, NOT contain 'ac' or '_')")

	upCmd := flag.NewFlagSet("up", flag.ExitOnError)

	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)

	if len(os.Args) < 2 {
		rootHelp()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "down":
		downCmd.Parse(os.Args[2:])

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

		// download
		dougaId := *downDougaId
		if dougaId == "" {
			downCmd.Usage()
			os.Exit(2)
		}

		config, err := config.LoadOrNewConfig()
		if err != nil {
			panic(err)
		}
		if err := config.MakeDownloadsHomeDir(); err != nil {
			panic(err)
		}

		if err := downloader.Download(config.DownloadsHomeDir, dougaId); err != nil {
			panic(err)
		}
	case "up":
		upCmd.Parse(os.Args[2:])
		panic("not implemented")
	case "version":
		versionCmd.Parse(os.Args[2:])
		fmt.Println(utils.GetVersion().Version)
	default:
		slog.Error("Invalid action")
		rootHelp()
		os.Exit(2)
	}
}
