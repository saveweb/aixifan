package utils

import "runtime/debug"

type Version struct {
	Version   string
	GoVersion string
}

// copy from github.com/internetarchive/Zeno (AGPLv3)
func GetVersion() (version Version) {
	version.Version = "unknown_version"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				version.Version = setting.Value
			}

			if setting.Key == "vcs.modified" {
				if setting.Value == "true" {
					version.Version += " (modified)"
				}
			}
		}
		version.GoVersion = info.GoVersion
	}
	return
}
