package downloader

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/canhlinh/hlsdl"
	"github.com/saveweb/aixifan/pkg/api"
	"github.com/saveweb/aixifan/pkg/extractor"
)

var headers = map[string]string{
	"User-Agent": "aixifanfan/0.0.1",
}

func validateDougaId(dougaId string) bool {
	if dougaId == "" {
		// return fmt.Errorf("dougaId is empty")
		slog.Error("dougaId is empty")
		return false
	}
	for _, c := range dougaId {
		if c < '0' || c > '9' {
			// return fmt.Errorf("dougaId should only contain digits")
			slog.Error("dougaId should only contain digits")
			return false
		}
	}
	return true
}

func SaveDougaInfos(dougaDir, dougaId string, parts []string) error {
	// save to dougaDir/ac{acid}.info.json
	for i, part := range parts {
		acid := fmt.Sprintf("%s_%d", dougaId, i+1)
		filepath := path.Join(dougaDir, "ac"+acid+".info.json")
		if err := os.WriteFile(filepath, []byte(part), 0644); err != nil {
			return err
		}
	}
	return nil
}

func DownloadVideo(dougaDir, acid, part string) error {
	dougaTitle, partTitle, err := extractor.GetTitles(part)
	if err != nil {
		return err
	}
	slog.Info("Downloading", "dougaTitle", dougaTitle, "partTitle", partTitle)

	ksPlayInfo, err := extractor.GetKsPlayJson(part)
	if err != nil {
		return err
	}

	m3u8s, err := extractor.GetM3U8s(ksPlayInfo)
	if err != nil {
		return err
	}
	if len(m3u8s) == 0 {
		return fmt.Errorf("m3u8s is empty")
	}

	m3u8 := m3u8s[0] // assume the first one is the best
	hlsDL := hlsdl.New(m3u8.Url, headers, dougaDir, "ac"+acid+".ts", 3, true)

	filepath, err := hlsDL.Download()
	if err != nil {
		return err
	}

	slog.Info("Downloaded", "filepath", filepath)
	return nil
}

func Download(downloadsHomeDir string, dougaId string) error {
	if !validateDougaId(dougaId) {
		return fmt.Errorf("invalid dougaId")
	}

	dougaDir := path.Join(downloadsHomeDir, dougaId)
	client := &http.Client{Timeout: 15 * time.Second}

	parts, err := api.GetDougaAll(client, dougaId)
	if err != nil {
		return err
	}
	if len(parts) == 0 {
		return fmt.Errorf("parts is empty")
	}

	if err := os.MkdirAll(dougaDir, 0755); err != nil {
		return err
	}

	if err := SaveDougaInfos(dougaDir, dougaId, parts); err != nil {
		slog.Error("SaveDougaInfos", "err", err)
		return err
	}
	slog.Info("DougaInfos saved", "dougaDir", dougaDir)

	for i, part := range parts {
		acid := fmt.Sprintf("%s_%d", dougaId, i+1)
		if err := DownloadVideo(dougaDir, acid, part); err != nil {
			return err
		}
	}
	return nil
}
