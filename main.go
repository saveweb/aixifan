package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/saveweb/aixifan/pkg/config"
	"github.com/saveweb/aixifan/pkg/downloader"
)

func help() {
	fmt.Println(
		`Usage: aixifan action dougaId

Action:
  down  Download
  up    Upload to IA
dougaId:
  int-string, NOT contain 'ac' or '_'

Examples:
e.g. aixifan down 32749`)
}

func main() {
	flag.Usage = help
	flag.Parse()
	if flag.NArg() < 2 {
		help()
		os.Exit(1)
	}

	action := flag.Arg(0)
	dougaId := flag.Arg(1)

	config, err := config.LoadOrNewConfig()
	if err != nil {
		panic(err)
	}
	if err := config.MakeDownloadsHomeDir(); err != nil {
		panic(err)
	}

	switch action {
	case "down":
		if err := downloader.Download(config.DownloadsHomeDir, dougaId); err != nil {
			panic(err)
		}
	case "up":
		panic("not implemented")
	default:
		panic("unknown action")
	}
}
