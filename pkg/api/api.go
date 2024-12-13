package api

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/saveweb/aixifan/pkg/extractor"
	"github.com/tidwall/gjson"
)

// dougaId: int-string
func requestDouga(client *http.Client, acId string) ([]byte, error) {
	if strings.Contains(acId, "ac") {
		return nil, errors.New("dougaId should not contain 'ac'")
	}
	url := "https://www.acfun.cn/v/ac" + acId + "?quickViewId=videoInfo_new&ajaxpipe=1"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "aixifanfan/0.0.1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// acId: 1234_1
func GetDouga(client *http.Client, acId string) (string, error) {
	slog.Info("GetDouga", "acId", acId)
	body, err := requestDouga(client, acId)
	if err != nil {
		return "", err
	}
	return extractor.Html2json(body)
}

// dougaId: 1234
func GetDougaAll(client *http.Client, dougaId string) ([]string, error) {
	slog.Info("GetDougaAll", "dougaId", dougaId)
	var parts []string
	if strings.Contains(dougaId, "_") {
		return nil, errors.New("dougaId should not contain '_'")
	}

	// get first part
	json, err := GetDouga(client, dougaId+"_1")
	if err != nil {
		return nil, err
	}
	parts = append(parts, json)
	// get len(videoList)
	videoCount := gjson.Get(json, "videoList.#").Int()
	for i := 2; i <= int(videoCount); i++ {
		// get next part
		json, err = GetDouga(client, dougaId+"_"+fmt.Sprint(i))
		if err != nil {
			return parts, err
		}
		parts = append(parts, json)
	}

	return parts, nil
}
