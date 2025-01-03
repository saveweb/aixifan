package downloader

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/canhlinh/hlsdl"
	"github.com/saveweb/aixifan/pkg/api"
	"github.com/saveweb/aixifan/pkg/extractor"
	"github.com/saveweb/aixifan/pkg/ffmpeg"
	"github.com/saveweb/aixifan/pkg/utils"
)

var headers = map[string]string{
	"User-Agent": utils.GetUA(),
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

// Cleanup hlsdl tmp files
func Cleanup(dougaDir string) error {
	infos, err := os.ReadDir(dougaDir)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if info.IsDir() {
			allDigits := true
			for _, c := range info.Name() {
				if c < '0' || c > '9' {
					allDigits = false
					break
				}
			}
			if allDigits {
				slog.Info("Removing cache", "dir", info.Name())
				if err := os.RemoveAll(path.Join(dougaDir, info.Name())); err != nil {
					return err
				}
			}
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

	qualityLabels := make([]string, len(m3u8s))
	for i, m3u8 := range m3u8s {
		qualityLabels[i] = m3u8.QualityLabel
	}
	slog.Info("Found m3u8s", "count", len(m3u8s), "QualityLabels", qualityLabels)
	m3u8 := m3u8s[0] // assume the first one is the best
	slog.Info("Selected", "QualityLabel", m3u8.QualityLabel)
	hlsDL := hlsdl.New(m3u8.Url, headers, dougaDir, "ac"+acid+".ts", 4, true)

	tsFilepath, err := hlsDL.Download()
	if err != nil {
		return err
	}

	slog.Info("Downloaded", "tsFilepath", tsFilepath)

	mp4Filepath, err := ffmpeg.TS2MP4(tsFilepath)
	if err != nil {
		slog.Error("TS2MP4", "err", err)
		return err
	}
	slog.Info("Converted", "mp4Filepath", mp4Filepath)

	if err := os.Remove(tsFilepath); err != nil {
		slog.Error("Remove tsFilepath", "err", err)
		return err
	}
	slog.Info("Removed", "tsFilepath", tsFilepath)

	return nil
}

func DownloadCover(dougaDir, part string) error {
	info, err := extractor.GetPartInfo(part)
	if err != nil {
		return err
	}
	urlParsed, err := url.Parse(info.CoverUrl)
	ext := path.Ext(urlParsed.Path)
	coverFilepath := path.Join(dougaDir, "douga"+info.DougaId+"_itemimage"+ext)
	slog.Info("Downloading cover", "coverFilepath", coverFilepath, "coverUrl", info.CoverUrl)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	req, _ := http.NewRequest("GET", info.CoverUrl, nil)
	req.Header.Set("User-Agent", utils.GetUA())

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := os.WriteFile(coverFilepath, body, 0644); err != nil {
		return err
	}

	return nil
}

func Download(downloadsHomeDir string, dougaId string) error {
	if !validateDougaId(dougaId) {
		return fmt.Errorf("invalid dougaId")
	}

	dougaDir := path.Join(downloadsHomeDir, dougaId)
	if isAllDownloaded(dougaDir) {
		slog.Info("Already downloaded", "dougaDir", dougaDir)
		return nil
	}

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

	if err := DownloadCover(dougaDir, parts[0]); err != nil {
		slog.Error("DownloadCover", "err", err)
		return err
	}
	slog.Info("DownloadCover saved", "dougaDir", dougaDir)

	for i, part := range parts {
		acid := fmt.Sprintf("%s_%d", dougaId, i+1)
		if err := DownloadVideo(dougaDir, acid, part); err != nil {
			return err
		}
	}
	if err := markAllDownloaded(dougaDir); err != nil {
		return err
	}

	return Cleanup(dougaDir)
}

// check _alldownloaded.mark
func isAllDownloaded(dougaDir string) bool {
	// check if dougaDir/_alldownloaded.mark exists
	_, err := os.Stat(path.Join(dougaDir, "_alldownloaded.mark"))
	return err == nil
}

func markAllDownloaded(dougaDir string) error {
	// create dougaDir/_alldownloaded.mark
	f, err := os.Create(path.Join(dougaDir, "_alldownloaded.mark"))
	if err != nil {
		return fmt.Errorf("failed to create _alldownloaded.mark: %w", err)
	}
	defer f.Close()
	return nil
}
