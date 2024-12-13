package ffmpeg

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

func IsAvailable() bool {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		slog.Error(err.Error())
		return false
	}
	return true
}
func TS2MP4(tsFilepath string) (string, error) {
	if !strings.HasSuffix(tsFilepath, ".ts") {
		return "", fmt.Errorf("tsFilepath should end with .ts")
	}

	mp4Filepath := tsFilepath[:len(tsFilepath)-2] + "mp4"
	if err := ts2mp4(tsFilepath, mp4Filepath); err != nil {
		return "", err
	}
	return mp4Filepath, nil
}

func ts2mp4(tsPath, mp4Path string) error {
	if !IsAvailable() {
		return fmt.Errorf("ffmpeg is not available")
	}

	tmpMp4Path := mp4Path + ".tmp"
	cmd := exec.Command("ffmpeg", "-y", "-i", tsPath, "-map", "0", "-c", "copy", "-f", "mp4", tmpMp4Path)
	if err := cmd.Run(); err != nil {
		slog.Error(err.Error())
		return err
	}
	os.Rename(tmpMp4Path, mp4Path)
	return nil
}
