package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type Version struct {
	Version   string // "unknown_version" or "<commit hash>" or "<commit hash> (modified)"
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

type GitTag struct {
	Name   string `json:"name"`
	Commit struct {
		Sha string `json:"sha"`
	} `json:"commit"`
}

func GetLatestTags() (tags []GitTag, err error) {
	client := &http.Client{Timeout: 3 * time.Second}
	for _, url := range []string{
		"https://git.saveweb.org/api/v1/repos/saveweb/aixifan/tags",
		"https://api.github.com/repos/saveweb/aixifan/tags", // in case git.saveweb.org is down
	} {
		var resp *http.Response
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "aixifan/"+GetVersion().Version)
		resp, err = client.Do(req)
		if err != nil {
			slog.Error(err.Error(), "url", url)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			slog.Error("Failed to get tags", "status", resp.Status, "url", url)
			continue
		}
		err = json.NewDecoder(resp.Body).Decode(&tags)
		if err != nil {
			slog.Error(err.Error(), "url", url)
			continue
		}
		return
	}
	return

}

func NewVersionAvailable() (bool, error) {
	tags, err := GetLatestTags()
	if err != nil {
		return false, err
	}
	latestTagName := tags[0].Name
	latestTagCommit := tags[0].Commit.Sha
	if latestTagCommit != GetVersion().Version {
		slog.Info("New version available", "latest_tag", latestTagName, "latest_commit", latestTagCommit)
		return true, nil
	}
	return false, nil
}
