package uploader

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/saveweb/aixifan/pkg/config"
	"github.com/saveweb/aixifan/pkg/extractor"
	"github.com/saveweb/aixifan/pkg/utils"
	iaupload "github.com/saveweb/go2internetarchive/pkg/upload"
	iautils "github.com/saveweb/go2internetarchive/pkg/utils"
)

func Main(upCmd *flag.FlagSet, upDougaId string, upDelete bool) int {
	dougaId := upDougaId
	if dougaId == "" {
		upCmd.Usage()
		return 2
	}

	config, err := config.LoadOrNewConfig()
	if err != nil {
		slog.Error("Failed to load config", "err", err)
		return 1
	}
	acckey, seckey, err := iautils.ReadKeysFromFile(config.IaKeyFile)
	if err != nil {
		slog.Error("Failed to read IA keys", "err", err)
		return 1
	}

	dougaDir := path.Join(config.DownloadsHomeDir, dougaId)

	dougaInfo, err := loadDougaInfoJson(dougaDir, dougaId, 1)
	if err != nil {
		slog.Error("Failed to load douga info", "err", err)
		return 1
	}

	videoCount := getVideoCount(dougaInfo)
	if videoCount == 0 {
		slog.Error("Video count is 0")
		return 1
	}

	for partNum := 1; partNum <= videoCount; partNum++ {
		if isUploadedMarked(dougaDir, dougaId, partNum) {
			slog.Info("Already uploaded, skipping...", "dougaId", dougaId, "partNum", partNum)
			continue
		}
		part, err := loadDougaInfoJson(dougaDir, dougaId, partNum)
		if err != nil {
			slog.Error("Failed to load douga info", "err", err)
			return 1
		}
		info, err := extractor.GetPartInfo(part)
		slog.Info("Uploading", "dougaId", dougaId, "partNum", partNum)
		if err != nil {
			slog.Error("Failed to get part info", "err", err)
			return 1
		}

		files, err := prepareFiles(dougaDir, dougaId, partNum)
		if err != nil {
			slog.Error("Failed to prepare files", "err", err)
			return 1
		}

		identifier := utils.ToIdentifier(dougaId, fmt.Sprint(partNum))

		externalIdentifiers := []string{
			fmt.Sprintf("urn:acfun:video:acid:%s", fmt.Sprintf("%s_%d", dougaId, partNum)),
			fmt.Sprintf("urn:acfun:video:dougaId:%s", dougaId),
			fmt.Sprintf("urn:acfun:video:videoId:%s", info.CurrentVideoInfo.Id),
			fmt.Sprintf("urn:acfun:video:userId:%s", info.User.Id),
		}

		subject := []string{
			"AcFun", "video",
		}
		for _, tag := range info.TagList {
			subject = append(subject, tag.Name)
		}

		// TODO: -collection
		collection := "opensource_movies"
		meta := map[string][]string{
			"mediatype":           {"movies"},
			"collection":          {collection},
			"title":               {fmt.Sprintf("%s P%d %s", info.DougaTitle, partNum, info.CurrentVideoInfo.PartTitle)},
			"description":         {info.Description},
			"creator":             {info.User.Name},
			"date":                {info.CreateTime},
			"external-identifier": externalIdentifiers,
			"subject":             subject,
			"originalurl":         {fmt.Sprintf("https://www.acfun.cn/v/ac%s_%d", dougaId, partNum)},
			"scanner":             {utils.GetUA()},
		}

		slog.Info("Uploading", "identifier", identifier, "meta", meta)

		err = iaupload.Upload(identifier, files, meta, acckey, seckey)
		if err != nil {
			slog.Error("Failed to upload", "err", err)
			return 1
		}

		slog.Info("Uploaded", "identifier", identifier, "archiveUrl", fmt.Sprintf("https://archive.org/details/%s", identifier))
		markAsUploaded(dougaDir, dougaId, partNum)
	}

	if upDelete {
		for countdown := 5; countdown > 0; countdown-- {
			fmt.Println()
			fmt.Printf("\rDeleting in %d seconds... (Ctrl+C to cancel)", countdown)
			time.Sleep(time.Second)
		}
		fmt.Println()

		if err := os.RemoveAll(dougaDir); err != nil {
			slog.Error("Failed to delete", "err", err)
			return 1
		}
		slog.Info("Deleted", "dougaDir", dougaDir)
	}
	return 1
}
