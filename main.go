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

func help() {
	fmt.Println("aixifan - " + utils.GetVersion().Version + " (" + utils.GetVersion().GoVersion + ")")
	fmt.Println(
		`Usage: aixifan action [dougaId]

Action:
  down  Download
  up    Upload to IA
  version  Show version
dougaId:
  int-string, NOT contain 'ac' or '_'

Examples:
e.g. aixifan down 32749`)
}

func main() {
	flag.Usage = help
	flag.Parse()

	action := flag.Arg(0)

	switch action {
	case "down":
		dougaId := flag.Arg(1)
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
		panic("not implemented")
	case "version":
		fmt.Println(utils.GetVersion().Version)
	default:
		slog.Error("Invalid action")
		help()
		os.Exit(2)
	}
}
