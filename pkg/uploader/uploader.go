package uploader

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/tidwall/gjson"
)

func markAsUploaded(dougaDir, dougaId string, partNum int) error {
	// ac{}_{}._uploaded.mark
	markFile := path.Join(dougaDir, fmt.Sprintf("ac%s_%d._uploaded.mark", dougaId, partNum))
	f, err := os.Create(markFile)
	if err != nil {
		return fmt.Errorf("failed to create mark file: %w", err)
	}
	defer f.Close()
	return nil
}

func isUploadedMarked(dougaDir, dougaId string, partNum int) bool {
	markFile := path.Join(dougaDir, fmt.Sprintf("ac%s_%d._uploaded.mark", dougaId, partNum))
	_, err := os.Stat(markFile)
	return err == nil
}

func loadDougaInfoJson(dougaDir, dougaId string, partNum int) (string, error) {
	acid_p1 := fmt.Sprintf("%s_%d", dougaId, partNum)

	f, err := os.Open(path.Join(dougaDir, "ac"+acid_p1+".info.json"))
	if err != nil {
		slog.Error("Failed to open info.json", "err", err)
		return "", err
	}
	defer f.Close()

	json, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read info.json: %w", err)
	}

	if !gjson.Valid(string(json)) {
		return "", fmt.Errorf("invalid json")
	}

	return string(json), nil
}

func getVideoCount(json string) int {
	return int(gjson.Get(json, "videoList.#").Int())
}

// ls dougaDir | grep douga{id}_itemimage*
func getCoverPath(dougaDir, dougaId string) string {
	coverpath := ""

	files, err := os.ReadDir(dougaDir)
	if err != nil {
		slog.Error("Failed to read dir", "err", err)
		return ""
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), fmt.Sprintf("douga%s_itemimage", dougaId)) {
			coverpath = path.Join(dougaDir, file.Name())
			break
		}
	}
	return coverpath
}

func prepareFiles(dougaDir, dougaId string, partNum int) (map[string]string, error) {
	files := map[string]string{}

	acid := fmt.Sprintf("%s_%d", dougaId, partNum)
	infopath := path.Join(dougaDir, "ac"+acid+".info.json")
	videopath := path.Join(dougaDir, "ac"+acid+".mp4")
	coverpath := getCoverPath(dougaDir, dougaId)

	for _, filepath := range []string{infopath, videopath, coverpath} {
		if filepath == "" {
			continue // coverpath is optional
		}
		if _, err := os.Stat(filepath); err != nil {
			slog.Error("File not found", "file", filepath)
			return nil, fmt.Errorf("file not found: %s", filepath)
		}
		remotePath := path.Base(filepath)
		files[remotePath] = filepath
	}

	return files, nil
}
