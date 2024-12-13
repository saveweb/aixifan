package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

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
	downSkipIACheck := downCmd.Bool("s", false, "Do not check if the item already exists on IA")

	upCmd := flag.NewFlagSet("up", flag.ExitOnError)

	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)

	if len(os.Args) < 2 {
		rootHelp()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "down":
		downCmd.Parse(os.Args[2:])
		os.Exit(downloader.Main(downCmd, downNoVersionCheck, downDougaId, downSkipIACheck))
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
